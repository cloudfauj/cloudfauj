package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/cloudfauj/cloudfauj/deployment"
	"github.com/cloudfauj/cloudfauj/environment"
	"github.com/cloudfauj/cloudfauj/infrastructure"
	"github.com/cloudfauj/cloudfauj/wsmanager"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

func (s *server) handlerDeployApp(w http.ResponseWriter, r *http.Request) {
	wsConn, err := s.wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log.Errorf("Failed to upgrade websocket connection: %v", err)
		return
	}
	defer wsConn.Close()
	conn := &wsmanager.WSManager{Conn: wsConn}

	var spec deployment.Spec
	if err := conn.ReadJSON(&spec); err != nil {
		s.log.Errorf("Failed to read deployment spec: %v", err)
		conn.SendFailureISE()
		return
	}
	if err := spec.CheckIsValid(); err != nil {
		conn.SendFailure(
			fmt.Sprintf("Invalid specification: %v", err),
			websocket.ClosePolicyViolation,
		)
		return
	}

	e, err := s.state.Environment(r.Context(), spec.TargetEnv)
	if err != nil {
		s.log.WithField("name", spec.TargetEnv).Errorf("Failed to check if target env exists: %v", err)
		conn.SendFailureISE()
		return
	}
	if e == nil {
		conn.SendFailure("Target environment does not exist", websocket.ClosePolicyViolation)
		return
	}
	if e.Status != environment.StatusProvisioned {
		conn.SendFailure(
			"Target environment is not ready to be deployed to",
			websocket.CloseInternalServerErr,
		)
		return
	}

	// create app dir inside env dir if it doesn't already exist
	dir := s.appTfDir(spec.TargetEnv, spec.App.Name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		s.log.Errorf("Failed to create directory for app: %v", err)
		conn.SendFailureISE()
		return
	}

	// create terraform object to run inside app directory
	tf, err := s.infra.NewTerraform(dir, conn)
	if err != nil {
		s.log.Error(err)
		conn.SendFailureISE()
		return
	}

	// get app from state if it already exists in the target environment
	app, err := s.state.App(r.Context(), spec.App.Name, spec.TargetEnv)
	if err != nil {
		s.log.WithField("name", spec.App.Name).Errorf("Failed to get app from state: %v", err)
		conn.SendFailureISE()
		return
	}
	if app == nil {
		s.createNewApp(r.Context(), conn, &spec, e, tf, dir)
		return
	}

	// app already exists, run new deployment
	s.deployApp(r.Context(), conn, &spec, tf)
}

func (s *server) createNewApp(
	ctx context.Context,
	conn *wsmanager.WSManager,
	spec *deployment.Spec,
	env *environment.Environment,
	tf *tfexec.Terraform,
	dir string,
) {
	s.log.WithFields(
		logrus.Fields{"name": spec.App.Name, "env": spec.TargetEnv},
	).Info("Creating new application")

	conn.SendTextMsg("Registering application in state")
	if err := s.state.CreateApp(ctx, spec.App, spec.TargetEnv); err != nil {
		s.log.Errorf("Failed to create app in state: %v", err)
		conn.SendFailureISE()
		return
	}

	i := &infrastructure.AppTFConfigInput{
		Spec:              spec,
		Env:               env,
		DomainTFStateFile: s.domainTFStateFile(env.Domain),
		EnvTFStateFile:    s.envTfStateFile(env.Name),
	}
	tfConfigs, err := s.infra.AppTFConfig(i)
	if err != nil {
		s.log.Errorf("Failed to generate terraform configurations for app: %v", err)
		conn.SendFailureISE()
		return
	}
	if err := s.writeFiles(dir, tfConfigs); err != nil {
		s.log.Errorf("Failed to write terraform configs for app: %v", err)
		conn.SendFailureISE()
		return
	}

	conn.SendTextMsg("Provisioning infrastructure")
	if err := s.infra.CreateApplication(ctx, spec, tf); err != nil {
		s.log.Errorf("Failed to provision app infrastructure: %v", err)
		conn.SendFailureISE()
		return
	}

	cluster, _ := s.infra.AppECSCluster(ctx, tf)
	service, _ := s.infra.AppECSService(ctx, tf)
	eventsCh := make(chan *Event)

	go s.trackDeployment(ctx, cluster, service, eventsCh)
	for e := range eventsCh {
		if e.Err != nil {
			conn.SendFailure(
				fmt.Sprintf("Deployment failed: %v", e.Err),
				websocket.CloseInternalServerErr,
			)
			return
		}
		conn.SendTextMsg(e.Msg)
	}

	conn.SendSuccess("App deployed successfully")
	return
}

func (s *server) deployApp(
	ctx context.Context,
	conn *wsmanager.WSManager,
	spec *deployment.Spec,
	tf *tfexec.Terraform,
) {
	depLogger := logrus.New()
	d := deployment.New(spec, depLogger)

	id, err := s.state.CreateDeployment(ctx, d)
	if err != nil {
		s.log.WithField("app", spec.App.Name).Errorf("Failed to create deployment: %v", err)
		conn.SendFailureISE()
		return
	}
	d.Id = strconv.FormatInt(id, 10)

	// open deployment log file
	if err := os.Mkdir(s.deploymentDir(d.Id), 0755); err != nil {
		s.log.WithField("deployment_id", d.Id).Errorf("Failed to create deployment dir: %v", err)
		conn.SendFailureISE()
		return
	}
	dlf, err := os.OpenFile(s.deploymentLogFile(d.Id), os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		// no need to exit, deployment logs can be written to default logger output
		s.log.Errorf("Failed to open deployment log file: %v", err)
	} else {
		defer dlf.Close()
		depLogger.SetOutput(dlf)
	}

	msg := "Deployment ID: " + d.Id
	conn.SendTextMsg(msg)
	d.Log(msg)
	s.log.WithFields(
		logrus.Fields{"name": spec.App.Name, "deployment_id": d.Id},
	).Info("Deploying application")

	if err := s.state.UpdateApp(ctx, spec.App, spec.TargetEnv); err != nil {
		s.log.Errorf("Failed to update app in state: %v", err)
		d.Fail(errors.New("a server error occurred while updating app state"))
		conn.SendFailureISE()
		return
	}
	if err := s.infra.ModifyApplication(ctx, spec, tf); err != nil {
		s.log.Errorf("Failed to modify application infrastructure: %v", err)
		conn.SendFailureISE()
		return
	}

	cluster, _ := s.infra.AppECSCluster(ctx, tf)
	service, _ := s.infra.AppECSService(ctx, tf)
	eventsCh := make(chan *Event)

	go s.trackDeployment(ctx, cluster, service, eventsCh)
	for e := range eventsCh {
		if e.Err != nil {
			d.Fail(e.Err)
			s.state.UpdateDeploymentStatus(ctx, d.Id, d.Status)
			conn.SendFailure(
				fmt.Sprintf("Deployment failed: %v", e.Err),
				websocket.CloseInternalServerErr,
			)
			return
		}
		d.Log(e.Msg)
		conn.SendTextMsg(e.Msg)
	}

	d.Succeed()
	s.state.UpdateDeploymentStatus(ctx, d.Id, d.Status)
	conn.SendSuccess("Deployed successfully")
}

func (s *server) handlerDestroyApp(w http.ResponseWriter, r *http.Request) {
	app := mux.Vars(r)["name"]
	env := r.URL.Query().Get("env")

	s.log.WithFields(
		logrus.Fields{"app": app, "env": env},
	).Info("App deletion request received")

	envState, err := s.state.Environment(r.Context(), env)
	if err != nil {
		s.log.Errorf("Failed to get environment from state: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if envState == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	appState, err := s.state.App(r.Context(), app, env)
	if err != nil {
		s.log.Errorf("Failed to get app from state: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if appState == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	appDir := s.appTfDir(env, app)

	// TODO: Use websocket in this controller and supply connection object below
	tf, err := s.infra.NewTerraform(appDir, nil)
	if err != nil {
		s.log.Errorf("Failed to create terraform object: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := s.infra.DestroyApplication(r.Context(), tf); err != nil {
		s.log.Errorf("Failed to destroy app infra: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := os.RemoveAll(appDir); err != nil {
		s.log.Errorf("Failed to delete app TF config from disk: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := s.state.DeleteApp(r.Context(), app, env); err != nil {
		s.log.Errorf("Failed to delete app from state: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.log.Info("Application deleted successfully")
	w.WriteHeader(http.StatusOK)
}

// trackDeployment polls the latest ECS deployment and streams the status until
// the deployment has completed or failed off.
func (s *server) trackDeployment(ctx context.Context, ecsCluster, ecsService string, e chan<- *Event) {
	defer close(e)

	// todo: improve timeout logic
	for j := 0; j < 120; j++ {
		e <- &Event{Msg: "Deploying application to ECS..."}
		d, err := s.infra.ECSServicePrimaryDeployment(ctx, ecsService, ecsCluster)
		if err != nil {
			s.log.Errorf("Failed to fetch deployment information from ECS: %v", err)
		}
		switch d.RolloutState {
		case types.DeploymentRolloutStateCompleted:
			e <- &Event{Msg: "Done"}
			return
		case types.DeploymentRolloutStateFailed:
			e <- &Event{Err: errors.New("ECS Deployment failed: " + aws.ToString(d.RolloutStateReason))}
			return
		}
		time.Sleep(time.Second * 5)
	}
	e <- &Event{Err: errors.New("deployment polling timeout reached")}
}

func (s *server) appTfDir(env, app string) string {
	return path.Join(s.envTfDir(env), app)
}

func (s *server) appTfFile(env, app string) string {
	return path.Join(s.appTfDir(env, app), s.config.terraformConfigFile)
}
