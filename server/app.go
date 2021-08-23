package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/cloudfauj/cloudfauj/deployment"
	"github.com/cloudfauj/cloudfauj/environment"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func (s *server) handlerDeployApp(w http.ResponseWriter, r *http.Request) {
	wsConn, err := s.wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log.Errorf("Failed to upgrade websocket connection: %v", err)
		return
	}
	defer wsConn.Close()
	conn := wsManager{wsConn}

	var spec deployment.Spec
	if err := conn.ReadJSON(&spec); err != nil {
		s.log.Errorf("Failed to read deployment spec: %v", err)
		conn.sendFailureISE()
		return
	}
	if err := spec.CheckIsValid(); err != nil {
		conn.sendFailure(
			fmt.Sprintf("Invalid specification: %v", err),
			websocket.ClosePolicyViolation,
		)
		return
	}

	e, err := s.state.Environment(r.Context(), spec.TargetEnv)
	if err != nil {
		s.log.WithField("name", spec.TargetEnv).Errorf("Failed to check if target env exists: %v", err)
		conn.sendFailureISE()
		return
	}
	if e == nil {
		conn.sendFailure("Target environment does not exist", websocket.ClosePolicyViolation)
		return
	}
	if e.Status != environment.StatusProvisioned {
		conn.sendFailure(
			"Target environment is not ready to be deployed to",
			websocket.CloseInternalServerErr,
		)
		return
	}

	// get app from state if it already exists in the target environment
	app, err := s.state.App(r.Context(), spec.App.Name, spec.TargetEnv)
	if err != nil {
		s.log.WithField("name", spec.App.Name).Errorf("Failed to get app from state: %v", err)
		conn.sendFailureISE()
		return
	}
	if app == nil {
		s.log.WithFields(
			logrus.Fields{"name": spec.App.Name, "env": spec.TargetEnv},
		).Info("Creating new application")

		conn.sendTextMsg("Creating application in this environment for the first time")

		if err := s.state.CreateApp(r.Context(), spec.App, spec.TargetEnv); err != nil {
			s.log.Errorf("Failed to create app in state: %v", err)
			conn.sendFailureISE()
			return
		}

		// provision infrastructure for app
		eventsCh := make(chan *Event)
		go s.provisionInfra(r.Context(), &spec, e, eventsCh)

		for e := range eventsCh {
			if e.Err != nil {
				conn.sendFailure(
					fmt.Sprintf("Creation failed: %v", e.Err),
					websocket.CloseInternalServerErr,
				)
				return
			}
			conn.sendTextMsg(e.Msg)
		}

		conn.sendSuccess("App created & deployed successfully")
		return
	}

	// app already exists, run new deployment

	depLogger := logrus.New()
	d := deployment.New(&spec, depLogger)

	id, err := s.state.CreateDeployment(r.Context(), d)
	if err != nil {
		s.log.WithField("app", spec.App.Name).Errorf("Failed to create deployment: %v", err)
		conn.sendFailureISE()
		return
	}
	d.Id = strconv.FormatInt(id, 10)

	// open deployment log file
	if err := os.Mkdir(s.deploymentDir(d.Id), 0755); err != nil {
		s.log.WithField("deployment_id", d.Id).Errorf("Failed to create deployment dir: %v", err)
		conn.sendFailureISE()
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
	conn.sendTextMsg(msg)
	d.Log(msg)
	s.log.
		WithFields(logrus.Fields{"name": spec.App.Name, "deployment_id": d.Id}).
		Info("Deploying application")

	if err := s.state.UpdateApp(r.Context(), spec.App, spec.TargetEnv); err != nil {
		s.log.Errorf("Failed to update app in state: %v", err)
		d.Fail(errors.New("a server error occurred while updating app state"))
		conn.sendFailureISE()
		return
	}

	// deploy app artifact
	eventsCh := make(chan *Event)
	go s.deployApp(r.Context(), &spec, e, eventsCh)

	for e := range eventsCh {
		if e.Err != nil {
			d.Fail(e.Err)
			s.state.UpdateDeploymentStatus(r.Context(), d.Id, d.Status)
			conn.sendFailure(
				fmt.Sprintf("Deployment failed: %v", e.Err),
				websocket.CloseInternalServerErr,
			)
			return
		}
		d.Log(e.Msg)
		conn.sendTextMsg(e.Msg)
	}

	d.Succeed()
	s.state.UpdateDeploymentStatus(r.Context(), d.Id, d.Status)
	conn.sendSuccess("Deployed successfully")
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

	if err := s.destroyInfra(r.Context(), env, app); err != nil {
		s.log.Errorf("Failed to destroy app infra: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := s.state.DeleteApp(r.Context(), app, env); err != nil {
		s.log.Errorf("Failed to delete app from state: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := os.Remove(s.appTfFile(env, app)); err != nil {
		s.log.Errorf("Failed to delete app TF config file from disk: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.log.Info("Application deleted successfully")
	w.WriteHeader(http.StatusOK)
}

func (s *server) provisionInfra(
	ctx context.Context,
	d *deployment.Spec,
	env *environment.Environment,
	e chan<- *Event,
) {
	defer close(e)

	e <- &Event{Msg: "provisioning infrastructure for application"}

	tf, err := s.infra.Tf(env.Name)
	if err != nil {
		e <- &Event{Err: err}
		return
	}
	tfConfig := s.infra.AppTfConfig(env.Name, d)
	if err := os.WriteFile(s.appTfFile(env.Name, d.App.Name), []byte(tfConfig), 0666); err != nil {
		e <- &Event{Err: fmt.Errorf("failed to create app terraform config : %v", err)}
		return
	}
	if err := tf.Init(ctx); err != nil {
		e <- &Event{Err: fmt.Errorf("failed to initialize terraform: %v", err)}
		return
	}
	e <- &Event{Msg: "Applying Terraform configuration"}
	if err := tf.Apply(ctx); err != nil {
		e <- &Event{Err: fmt.Errorf("failed to apply terraform changes: %v", err)}
		return
	}

	res, err := tf.Output(ctx)
	if err != nil {
		e <- &Event{Err: fmt.Errorf("failed to read terraform output: %v", err)}
		return
	}
	cluster := string(res["ecs_cluster_arn"].Value)
	cluster = strings.Trim(cluster, "\"")
	service := string(res[d.App.Name+"_ecs_service"].Value)
	service = strings.Trim(service, "\"")

	s.trackDeployment(ctx, cluster, service, e)
}

func (s *server) deployApp(
	ctx context.Context,
	d *deployment.Spec,
	env *environment.Environment,
	e chan<- *Event,
) {
	defer close(e)

	tf, err := s.infra.Tf(env.Name)
	if err != nil {
		e <- &Event{Err: err}
		return
	}
	// overwrite the existing app tf config with new one
	tfConfig := s.infra.AppTfConfig(env.Name, d)
	if err := os.WriteFile(s.appTfFile(env.Name, d.App.Name), []byte(tfConfig), 0666); err != nil {
		e <- &Event{Err: fmt.Errorf("failed to create app terraform config : %v", err)}
		return
	}
	e <- &Event{Msg: "Applying Terraform configuration"}
	if err := tf.Apply(ctx); err != nil {
		e <- &Event{Err: fmt.Errorf("failed to apply terraform changes: %v", err)}
		return
	}
	res, err := tf.Output(ctx)
	if err != nil {
		e <- &Event{Err: fmt.Errorf("failed to read terraform output: %v", err)}
		return
	}
	cluster := string(res["ecs_cluster_arn"].Value)
	cluster = strings.Trim(cluster, "\"")
	service := string(res[d.App.Name+"_ecs_service"].Value)
	service = strings.Trim(service, "\"")

	s.trackDeployment(ctx, cluster, service, e)
}

// trackDeployment polls the latest ECS deployment and streams the status until
// the deployment has completed or failed off.
func (s *server) trackDeployment(ctx context.Context, ecsCluster, ecsService string, e chan<- *Event) {
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

func (s *server) destroyInfra(ctx context.Context, env, app string) error {
	tf, err := s.infra.Tf(env)
	if err != nil {
		return err
	}
	if err := tf.Destroy(ctx, tfexec.Target("module."+app)); err != nil {
		return fmt.Errorf("failed to destroy terraform infra: %v", err)
	}
	return nil
}

func (s *server) appTfFile(env, app string) string {
	return path.Join(s.envTfDir(env), app+".tf")
}
