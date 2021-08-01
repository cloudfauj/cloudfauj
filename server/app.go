package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudfauj/cloudfauj/application"
	"github.com/cloudfauj/cloudfauj/deployment"
	"github.com/cloudfauj/cloudfauj/environment"
	infra "github.com/cloudfauj/cloudfauj/infrastructure"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
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
		s.log.Errorf("Invalid specification: %v", err)
		_ = sendWSClosureMsg(conn, websocket.CloseInternalServerErr)
		return
	}

	e, err := s.state.Environment(r.Context(), spec.TargetEnv)
	if err != nil {
		s.log.WithField("name", spec.TargetEnv).Errorf("Failed to check if target env exists: %v", err)
		_ = sendWSClosureMsg(conn, websocket.CloseInternalServerErr)
		return
	}
	if e == nil {
		s.log.WithField("name", spec.TargetEnv).Debug("Deployment target environment does not exist")
		_ = sendWSClosureMsg(conn, websocket.ClosePolicyViolation)
		return
	}
	if e.Status != environment.StatusProvisioned {
		s.log.WithField("name", spec.TargetEnv).Error("Environment is not ready to be deployed to")
		_ = sendWSClosureMsg(conn, websocket.CloseInternalServerErr)
		return
	}

	// get app from state if it already exists
	app, err := s.state.App(r.Context(), spec.App.Name)
	if err != nil {
		s.log.WithField("name", spec.App.Name).Errorf("Failed to get app from state: %v", err)
		_ = sendWSClosureMsg(conn, websocket.CloseInternalServerErr)
		return
	}
	if app == nil {
		s.log.WithField("name", spec.App.Name).Info("Creating new application")

		// register app in state
		if err := s.state.CreateApp(r.Context(), spec.App); err != nil {
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

		_ = sendWSClosureMsg(conn, websocket.CloseNormalClosure)
		return
	}

	// create new deployment in state
	d := deployment.New(&spec)
	id, err := s.state.CreateDeployment(r.Context(), d)
	if err != nil {
		s.log.WithField("app", spec.App.Name).Errorf("Failed to create deployment: %v", err)
		_ = sendWSClosureMsg(conn, websocket.CloseInternalServerErr)
		return
	}
	d.Id = strconv.FormatInt(id, 10)

	// todo: log ID in dep log & send to websocket
	s.log.
		WithFields(logrus.Fields{"name": spec.App.Name, "deployment_id": d.Id}).
		Info("Deploying application")

	if err := s.state.UpdateApp(r.Context(), spec.App); err != nil {
		s.log.Errorf("Failed to update app in state: %v", err)
		_ = sendWSClosureMsg(conn, websocket.CloseInternalServerErr)
		return
	}

	// provision infrastructure for app
	eventsCh := make(chan *Event)
	go s.deployApp(r.Context(), &spec, app, e, eventsCh)

	for e := range eventsCh {
		if e.Err != nil {
			d.Fail(e.Err)
			s.state.UpdateDeploymentStatus(r.Context(), d.Id, d.Status)

			m := []byte(fmt.Sprintf("Deployment failed: %v", e.Err))
			conn.WriteMessage(websocket.TextMessage, m)

			return
		}
		d.AppendLog(e.Msg)
		conn.WriteMessage(websocket.TextMessage, []byte(e.Msg))
	}

	s.state.UpdateDeploymentStatus(r.Context(), d.Id, deployment.StatusSucceeded)
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

	appState, err := s.state.App(r.Context(), app)
	if err != nil {
		s.log.Errorf("Failed to get app from state: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if appState == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	appInfra, err := s.state.AppInfra(r.Context(), app)
	if err != nil {
		s.log.Errorf("Failed to get app infra from state: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := s.destroyInfra(r.Context(), appInfra, envState); err != nil {
		s.log.Errorf("Failed to destroy app infra: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := s.state.DeleteAppInfra(r.Context(), app); err != nil {
		s.log.Errorf("Failed to delete app infra from state: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := s.state.DeleteApp(r.Context(), app); err != nil {
		s.log.Errorf("Failed to delete app from state: %v", err)
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

	i := &infra.AppInfra{App: d.App.Name}
	e <- &Event{Msg: "provisioning infrastructure for public application"}

	td, err := s.infra.CreateTaskDefinition(ctx, &infra.TaskDefintionParams{
		Env:          env.Name,
		Service:      d.App.Name,
		TaskExecRole: env.Res.TaskExecIAMRole,
		Image:        d.Artifact,
		Cpu:          d.App.Resources.Cpu,
		Memory:       d.App.Resources.Memory,
		BindPort:     d.App.Resources.Network.BindPort,
	})
	if err != nil {
		e <- &Event{Err: fmt.Errorf("failed to create task definition: %v", err)}
		return
	}
	i.EcsTaskDefinition = td
	e <- &Event{Msg: "created ECS task definition"}

	sg, err := s.infra.CreateSecurityGroup(ctx, env.Name, d.App.Name, env.Res.VpcId, d.App.Resources.Network.BindPort)
	if err != nil {
		e <- &Event{Err: fmt.Errorf("failed to create security group: %v", err)}
		return
	}
	i.SecurityGroup = sg
	e <- &Event{Msg: "created security group"}

	srv, err := s.infra.CreateECSService(ctx, &infra.ECSServiceParams{
		Env:           env.Name,
		Service:       d.App.Name,
		Cluster:       env.Res.ECSCluster,
		TaskDef:       i.EcsTaskDefinition,
		ComputeSubnet: env.Res.ComputeSubnet,
		SecurityGroup: i.SecurityGroup,
	})
	if err != nil {
		e <- &Event{
			Err: fmt.Errorf("failed to create ECS service: %v", err),
		}
		return
	}
	i.ECSService = srv
	e <- &Event{Msg: "created ECS service"}

	// todo: do we need to setup autoscaling separately? with fargate?

	// register all infra resources in state
	if err := s.state.CreateAppInfra(ctx, i); err != nil {
		e <- &Event{Err: fmt.Errorf("failed to register infra in state: %v", err)}
		return
	}

	// todo: tail ecs deployment logs
}

func (s *server) deployApp(
	ctx context.Context,
	d *deployment.Spec,
	originalApp *application.Application,
	env *environment.Environment,
	e chan<- *Event,
) {
	defer close(e)

	if d.App.Resources.Network.BindPort != originalApp.Resources.Network.BindPort {
		e <- &Event{Err: errors.New("changing bind port of application is not supported")}
		return
	}

	i, err := s.state.AppInfra(ctx, d.App.Name)
	if err != nil {
		e <- &Event{Err: fmt.Errorf("failed to fetch app state: %v", err)}
		return
	}

	td, err := s.infra.CreateTaskDefinition(ctx, &infra.TaskDefintionParams{
		Env:          env.Name,
		Service:      d.App.Name,
		TaskExecRole: env.Res.TaskExecIAMRole,
		Image:        d.Artifact,
		Cpu:          d.App.Resources.Cpu,
		Memory:       d.App.Resources.Memory,
		BindPort:     d.App.Resources.Network.BindPort,
	})
	if err != nil {
		e <- &Event{Err: fmt.Errorf("failed to create new task definition: %v", err)}
		return
	}
	i.EcsTaskDefinition = td
	e <- &Event{Msg: "created new ECS task definition"}

	if err := s.infra.UpdateECSService(ctx, d.App.Name, env.Res.ECSCluster, i.EcsTaskDefinition); err != nil {
		e <- &Event{Err: fmt.Errorf("failed to update ECS service: %v", err)}
		return
	}
	e <- &Event{Msg: "updated ECS service"}

	if err := s.state.UpdateAppInfra(ctx, i); err != nil {
		e <- &Event{Err: fmt.Errorf("failed to update app infra state: %v", err)}
		return
	}
	e <- &Event{Msg: "updated app infra state"}

	// todo: tail deployment logs

	e <- &Event{Msg: "deployment succeeded"}
}

func (s *server) destroyInfra(ctx context.Context, app *infra.AppInfra, env *environment.Environment) error {
	// todo: determine if we want to delete task definition(s) as well
	//  this has been left out for now

	if err := s.infra.DrainECSService(ctx, app.ECSService, env.Res.ECSCluster); err != nil {
		return fmt.Errorf("failed to drain ECS Service: %v", err)
	}
	// todo: improve the timeout logic
	for i := 0; i < 12; i++ {
		s.log.Info("Draining ECS service...")
		time.Sleep(time.Second * 5)
		c, err := s.infra.ECSServiceTaskCount(ctx, app.ECSService, env.Res.ECSCluster)
		if err == nil && c == 0 {
			s.log.Info("Done")
			break
		}
	}
	if err := s.infra.DestroyECSService(ctx, app.ECSService, env.Res.ECSCluster); err != nil {
		return fmt.Errorf("failed to delete ECS Service: %v", err)
	}
	if err := s.infra.DestroySecurityGroup(ctx, app.SecurityGroup); err != nil {
		return fmt.Errorf("failed to delete security group: %v", err)
	}
	return nil
}
