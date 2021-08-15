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
	"time"
)

func (s *server) handlerDeployApp(w http.ResponseWriter, r *http.Request) {
	conn, _ := s.wsUpgrader.Upgrade(w, r, nil)
	defer conn.Close()

	var spec deployment.Spec
	if err := conn.ReadJSON(&spec); err != nil {
		s.log.Errorf("Failed to read deployment spec: %v", err)
		_ = sendWSClosureMsg(conn, websocket.CloseInternalServerErr)
		return
	}
	if err := spec.CheckIsValid(); err != nil {
		conn.WriteMessage(
			websocket.TextMessage,
			[]byte(fmt.Sprintf("Invalid specification: %v", err)),
		)
		_ = sendWSClosureMsg(conn, websocket.ClosePolicyViolation)
		return
	}

	e, err := s.state.Environment(r.Context(), spec.TargetEnv)
	if err != nil {
		s.log.WithField("name", spec.TargetEnv).Errorf("Failed to check if target env exists: %v", err)
		_ = sendWSClosureMsg(conn, websocket.CloseInternalServerErr)
		return
	}
	if e == nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Target environment does not exist"))
		_ = sendWSClosureMsg(conn, websocket.ClosePolicyViolation)
		return
	}
	if e.Status != environment.StatusProvisioned {
		conn.WriteMessage(websocket.TextMessage, []byte("Target environment is not ready to be deployed to"))
		_ = sendWSClosureMsg(conn, websocket.CloseInternalServerErr)
		return
	}

	// get app from state if it already exists in the target environment
	app, err := s.state.App(r.Context(), spec.App.Name, spec.TargetEnv)
	if err != nil {
		s.log.WithField("name", spec.App.Name).Errorf("Failed to get app from state: %v", err)
		_ = sendWSClosureMsg(conn, websocket.CloseInternalServerErr)
		return
	}
	if app == nil {
		s.log.WithFields(
			logrus.Fields{"name": spec.App.Name, "env": spec.TargetEnv},
		).Info("Creating new application")
		conn.WriteMessage(
			websocket.TextMessage,
			[]byte("Creating application in this environment for the first time"),
		)

		// register app in state
		if err := s.state.CreateApp(r.Context(), spec.App, spec.TargetEnv); err != nil {
			s.log.Errorf("Failed to create app in state: %v", err)
			_ = sendWSClosureMsg(conn, websocket.CloseInternalServerErr)
			return
		}

		// provision infrastructure for app
		eventsCh := make(chan *Event)
		go s.provisionInfra(r.Context(), &spec, e, eventsCh)

		for e := range eventsCh {
			if e.Err != nil {
				m := []byte(fmt.Sprintf("Creation failed: %v", e.Err))
				conn.WriteMessage(websocket.TextMessage, m)
				return
			}
			conn.WriteMessage(websocket.TextMessage, []byte(e.Msg))
		}

		conn.WriteMessage(websocket.TextMessage, []byte("App created & deployed successfully"))
		_ = sendWSClosureMsg(conn, websocket.CloseNormalClosure)
		return
	}

	// app already exists, run new deployment

	depLogger := logrus.New()
	d := deployment.New(&spec, depLogger)

	id, err := s.state.CreateDeployment(r.Context(), d)
	if err != nil {
		s.log.WithField("app", spec.App.Name).Errorf("Failed to create deployment: %v", err)
		_ = sendWSClosureMsg(conn, websocket.CloseInternalServerErr)
		return
	}
	d.Id = strconv.FormatInt(id, 10)

	// open deployment log file
	os.Mkdir(s.deploymentDir(d.Id), 0755)
	dlf, err := os.OpenFile(s.deploymentLogFile(d.Id), os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		s.log.Errorf("Failed to open deployment log file: %v", err)
	} else {
		defer dlf.Close()
		depLogger.SetOutput(dlf)
	}

	msg := "Deployment ID: " + d.Id
	conn.WriteMessage(websocket.TextMessage, []byte(msg))
	d.Log(msg)
	s.log.
		WithFields(logrus.Fields{"name": spec.App.Name, "deployment_id": d.Id}).
		Info("Deploying application")

	if err := s.state.UpdateApp(r.Context(), spec.App, spec.TargetEnv); err != nil {
		s.log.Errorf("Failed to update app in state: %v", err)
		d.Fail(errors.New("a server error occurred while updating app state"))
		_ = sendWSClosureMsg(conn, websocket.CloseInternalServerErr)
		return
	}

	// deploy app artifact
	eventsCh := make(chan *Event)
	go s.deployApp(r.Context(), &spec, e, eventsCh)

	for e := range eventsCh {
		if e.Err != nil {
			d.Fail(e.Err)
			s.state.UpdateDeploymentStatus(r.Context(), d.Id, d.Status)

			m := []byte(fmt.Sprintf("Deployment failed: %v", e.Err))
			conn.WriteMessage(websocket.TextMessage, m)

			return
		}
		d.Log(e.Msg)
		conn.WriteMessage(websocket.TextMessage, []byte(e.Msg))
	}

	d.Succeed()
	s.state.UpdateDeploymentStatus(r.Context(), d.Id, d.Status)

	conn.WriteMessage(websocket.TextMessage, []byte("Deployed successfully"))
	_ = sendWSClosureMsg(conn, websocket.CloseNormalClosure)
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

	tf := s.infra.AppTfConfig(env.Name, d)
	if err := os.WriteFile(s.appTfFile(env.Name, d.App.Name), []byte(tf), 0666); err != nil {
		e <- &Event{Err: fmt.Errorf("failed to create app terraform config : %v", err)}
		return
	}
	if err := s.infra.Tf.Init(ctx); err != nil {
		e <- &Event{Err: fmt.Errorf("failed to initialize terraform: %v", err)}
		return
	}
	e <- &Event{Msg: "Applying Terraform configuration"}
	module := fmt.Sprintf("module.%s_%s", env.Name, d.App.Name)
	if err := s.infra.Tf.Apply(ctx, tfexec.Target(module)); err != nil {
		e <- &Event{Err: fmt.Errorf("failed to apply terraform changes: %v", err)}
		return
	}

	res, err := s.infra.Tf.Output(ctx)
	if err != nil {
		e <- &Event{Err: fmt.Errorf("failed to read terraform output: %v", err)}
		return
	}
	cluster, _ := res[fmt.Sprintf("%s_ecs_cluster_arn", env.Name)].Value.MarshalJSON()
	service, _ := res[fmt.Sprintf("%s_%s_ecs_service", env.Name, d.App.Name)].Value.MarshalJSON()

	s.trackDeployment(ctx, string(cluster), string(service), e)
}

func (s *server) deployApp(
	ctx context.Context,
	d *deployment.Spec,
	env *environment.Environment,
	e chan<- *Event,
) {
	defer close(e)

	// overwrite the existing app tf config with new one
	tf := s.infra.AppTfConfig(env.Name, d)
	if err := os.WriteFile(s.appTfFile(env.Name, d.App.Name), []byte(tf), 0666); err != nil {
		e <- &Event{Err: fmt.Errorf("failed to create app terraform config : %v", err)}
		return
	}
	e <- &Event{Msg: "Applying Terraform configuration"}
	module := fmt.Sprintf("module.%s_%s", env.Name, d.App.Name)
	if err := s.infra.Tf.Apply(ctx, tfexec.Target(module)); err != nil {
		e <- &Event{Err: fmt.Errorf("failed to apply terraform changes: %v", err)}
		return
	}
	res, err := s.infra.Tf.Output(ctx)
	if err != nil {
		e <- &Event{Err: fmt.Errorf("failed to read terraform output: %v", err)}
		return
	}
	cluster := res[fmt.Sprintf("%s_ecs_cluster_arn", env.Name)].Value
	service := res[fmt.Sprintf("%s_%s_ecs_service", env.Name, d.App.Name)].Value

	s.trackDeployment(ctx, string(cluster), string(service), e)
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
	module := fmt.Sprintf("module.%s_%s", env, app)
	if err := s.infra.Tf.Destroy(ctx, tfexec.Target(module)); err != nil {
		return fmt.Errorf("failed to destroy terraform infra: %v", err)
	}
	return nil
}

func (s *server) appTfFile(env, app string) string {
	return path.Join(s.config.DataDir, TerraformDir, fmt.Sprintf("%s_%s.tf", env, app))
}
