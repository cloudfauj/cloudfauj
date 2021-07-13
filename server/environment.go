package server

import (
	"github.com/cloudfauj/cloudfauj/environment"
	"github.com/gorilla/websocket"
	"net/http"
)

func (s *server) handlerListEnvironments(w http.ResponseWriter, r *http.Request) {}

func (s *server) handlerCreateEnv(w http.ResponseWriter, r *http.Request) {
	conn, _ := s.wsUpgrader.Upgrade(w, r, nil)
	defer conn.Close()

	var env environment.Environment
	if err := conn.ReadJSON(&env); err != nil {
		s.log.Errorf("Failed to read env config: %v", err)
		_ = sendWSClosureMsg(conn, websocket.CloseInternalServerErr)
		return
	}
	if err := env.CheckIsValid(); err != nil {
		s.log.Infof("Invalid env config: %v", err)
		_ = sendWSClosureMsg(conn, websocket.CloseInvalidFramePayloadData)
		return
	}

	ok, err := s.state.CheckEnvExists(r.Context(), env.Name)
	if err != nil {
		s.log.Errorf("Failed to check if env exists: %v", err)
		_ = sendWSClosureMsg(conn, websocket.CloseInternalServerErr)
		return
	}
	if ok {
		s.log.WithField("name", env.Name).Debug("Environment already exists")
		_ = sendWSClosureMsg(conn, websocket.ClosePolicyViolation)
		return
	}

	s.log.WithField("name", env.Name).Info("Creating new environment")

	env.Status = environment.StatusProvisioning
	if err := s.state.CreateEnvironment(r.Context(), &env); err != nil {
		s.log.Errorf("Failed to store env info in state: %v", err)
		_ = sendWSClosureMsg(conn, websocket.CloseInternalServerErr)
		return
	}

	eventsCh := make(chan environment.Event)
	go env.Provision(r.Context(), eventsCh)

	for e := range eventsCh {
		if e.Err != nil {
			s.log.Errorf("Failed to provision environment: %v", e.Err)
			_ = sendWSClosureMsg(conn, websocket.CloseInternalServerErr)
			return
		}
		s.log.Info(e.Msg)
		_ = conn.WriteMessage(websocket.TextMessage, []byte(e.Msg))
	}

	env.Status = environment.StatusProvisioned
	if err := s.state.UpdateEnvironment(r.Context(), &env); err != nil {
		s.log.Errorf("Failed to update env info in state: %v", err)
		_ = sendWSClosureMsg(conn, websocket.CloseInternalServerErr)
		return
	}

	_ = sendWSClosureMsg(conn, websocket.CloseNormalClosure)
}

func (s *server) handlerDestroyEnv(w http.ResponseWriter, r *http.Request) {}
