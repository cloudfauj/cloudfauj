package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudfauj/cloudfauj/application"
	"github.com/cloudfauj/cloudfauj/deployment"
	infra "github.com/cloudfauj/cloudfauj/infrastructure"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"io/fs"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
)

func (s *server) handlerGetAppLogs(w http.ResponseWriter, r *http.Request) {
	app := mux.Vars(r)["name"]
	env := r.URL.Query().Get("env")
	f := path.Join(s.config.DataDir, ApplicationsDir, app, ApplicationsEnvDir, env, LogFileBasename)

	s.log.WithField("path", f).Info("Fetching application logs")

	content, err := ioutil.ReadFile(f)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		s.log.Errorf("Failed to read log file: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := strings.Split(string(content), "\n")
	jsonRes, _ := json.Marshal(res)
	_, _ = w.Write(jsonRes)
}

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

	// ensure target environment exists
	ok, err := s.state.CheckEnvExists(r.Context(), spec.TargetEnv)
	if err != nil {
		s.log.WithField("name", spec.TargetEnv).Errorf("Failed to check if target env exists: %v", err)
		_ = sendWSClosureMsg(conn, websocket.CloseInternalServerErr)
		return
	}
	if !ok {
		s.log.WithField("name", spec.TargetEnv).Debug("Deployment target environment does not exist")
		_ = sendWSClosureMsg(conn, websocket.ClosePolicyViolation)
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
	d.Id = id
	// todo: log ID in dep log & send to websocket

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
		go s.provisionInfra(r.Context(), spec.App, eventsCh)

		for e := range eventsCh {
			if e.Err != nil {
				d.Fail(e.Err)
				s.state.UpdateDeploymentStatus(r.Context(), d.Status)

				m := []byte(fmt.Sprintf("Deployment failed: %v", e.Err))
				conn.WriteMessage(websocket.TextMessage, m)

				return
			}
			d.AppendLog(e.Msg)
			conn.WriteMessage(websocket.TextMessage, []byte(e.Msg))
		}
	}

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
	go s.deployApp(r.Context(), &spec, eventsCh)

	for e := range eventsCh {
		if e.Err != nil {
			d.Fail(e.Err)
			s.state.UpdateDeploymentStatus(r.Context(), d.Status)

			m := []byte(fmt.Sprintf("Deployment failed: %v", e.Err))
			conn.WriteMessage(websocket.TextMessage, m)

			return
		}
		d.AppendLog(e.Msg)
		conn.WriteMessage(websocket.TextMessage, []byte(e.Msg))
	}

	s.state.UpdateDeploymentStatus(r.Context(), deployment.StatusSucceeded)
	_ = sendWSClosureMsg(conn, websocket.CloseNormalClosure)
}

func (s *server) provisionInfra(ctx context.Context, a *application.Application, e chan<- *Event) {
	defer close(e)

	i := &infra.AppInfra{App: a.Name}

	// ensure public visibility
	if a.Visibility != application.VisibilityPublic {
		e <- &Event{Err: errors.New("only public visibility is supported for app")}
		return
	}
	e <- &Event{Msg: "provisioning infrastructure for public application"}

	// create task definition
	td, err := s.infra.CreateTaskDefinition(ctx)
	if err != nil {
		e <- &Event{Err: fmt.Errorf("failed to create task definition: %v", err)}
		return
	}
	i.EcsTaskDefinition = td
	e <- &Event{Msg: "created ECS task definition"}

	// create target group
	tg, err := s.infra.CreateTargetGroup(ctx)
	if err != nil {
		e <- &Event{Err: fmt.Errorf("failed to create target group: %v", err)}
		return
	}
	i.TargetGroup = tg
	e <- &Event{Msg: "created target group"}

	// attach target group to ALB
	rule, err := s.infra.AttachTargetGroup(ctx, tg)
	if err != nil {
		e <- &Event{
			Err: fmt.Errorf("failed to attach target group to load balancer: %v", err),
		}
		return
	}
	i.AlbListenerRule = rule
	e <- &Event{Msg: "attached target group to load balancer"}

	// create DNS record
	rec, err := s.infra.CreateDNSRecord(ctx)
	if err != nil {
		e <- &Event{
			Err: fmt.Errorf("failed to create DNS record: %v", err),
		}
		return
	}
	i.DNSRecord = rec
	e <- &Event{Msg: "created DNS record"}

	// create ECS service
	srv, err := s.infra.CreateECSService(ctx)
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

func (s *server) deployApp(ctx context.Context, spec *deployment.Spec, e chan<- *Event) {
	defer close(e)

	// ensure public visibility
	if spec.App.Visibility != application.VisibilityPublic {
		e <- &Event{Err: errors.New("only public visibility is supported for app")}
		return
	}

	i, err := s.state.AppInfra(ctx, spec.App.Name)
	if err != nil {
		e <- &Event{Err: fmt.Errorf("failed to fetch app state: %v", err)}
		return
	}

	// create task definition
	td, err := s.infra.CreateTaskDefinition(ctx)
	if err != nil {
		e <- &Event{Err: fmt.Errorf("failed to create new task definition: %v", err)}
		return
	}
	i.EcsTaskDefinition = td
	e <- &Event{Msg: "created new ECS task definition"}

	// update ecs service
	if err := s.infra.UpdateECSService(ctx, td); err != nil {
		e <- &Event{Err: fmt.Errorf("failed to update ECS service: %v", err)}
		return
	}
	e <- &Event{Msg: "updated ECS service"}

	// update app state with new taskdef
	if err := s.state.UpdateAppInfra(ctx, i); err != nil {
		e <- &Event{Err: fmt.Errorf("failed to update app infra state: %v", err)}
		return
	}
	e <- &Event{Msg: "updated app infra state"}

	// todo: tail deployment logs

	e <- &Event{Msg: "deployment succeeded"}
}
