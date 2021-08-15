package server

import (
	"encoding/json"
	"fmt"
	"github.com/cloudfauj/cloudfauj/environment"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
	"path"
)

func (s *server) handlerListEnvironments(w http.ResponseWriter, r *http.Request) {
	res, err := s.state.ListEnvironments(r.Context())
	if err != nil {
		s.log.Errorf("Failed to list environments from state: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsonRes, _ := json.Marshal(res)
	_, _ = w.Write(jsonRes)
}

func (s *server) handlerCreateEnv(w http.ResponseWriter, r *http.Request) {
	var env *environment.Environment

	conn, _ := s.wsUpgrader.Upgrade(w, r, nil)
	defer conn.Close()

	if err := conn.ReadJSON(&env); err != nil {
		s.log.Errorf("Failed to read env config: %v", err)
		_ = sendWSClosureMsg(conn, websocket.CloseInternalServerErr)
		return
	}
	if err := env.CheckIsValid(); err != nil {
		conn.WriteMessage(
			websocket.TextMessage,
			[]byte(fmt.Sprintf("Invalid environment config: %v", err)),
		)
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
		conn.WriteMessage(websocket.TextMessage, []byte("Environment already exists"))
		_ = sendWSClosureMsg(conn, websocket.ClosePolicyViolation)
		return
	}

	s.log.WithField("name", env.Name).Info("Creating new environment")

	env.Status = environment.StatusProvisioning
	env.Infra = s.infra

	if err := s.state.CreateEnvironment(r.Context(), env); err != nil {
		s.log.Errorf("Failed to store env info in state: %v", err)
		_ = sendWSClosureMsg(conn, websocket.CloseInternalServerErr)
		return
	}
	conn.WriteMessage(websocket.TextMessage, []byte("Registered in state"))

	f, err := os.OpenFile(s.envTfFile(env.Name), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		s.log.Errorf("Failed to create Terraform config file for env: %v", err)
		_ = sendWSClosureMsg(conn, websocket.CloseInternalServerErr)
		return
	}
	defer f.Close()

	eventsCh := make(chan environment.Event)
	go env.Provision(r.Context(), f, eventsCh)

	for e := range eventsCh {
		if e.Err != nil {
			s.log.Errorf("Failed to provision environment: %v", e.Err)
			_ = sendWSClosureMsg(conn, websocket.CloseInternalServerErr)
			return
		}
		s.log.Info(e.Msg)
		conn.WriteMessage(websocket.TextMessage, []byte(e.Msg))
	}

	env.Status = environment.StatusProvisioned
	if err := s.state.UpdateEnvironment(r.Context(), env); err != nil {
		s.log.Errorf("Failed to update env info in state: %v", err)
		_ = sendWSClosureMsg(conn, websocket.CloseInternalServerErr)
		return
	}

	conn.WriteMessage(websocket.TextMessage, []byte("Successfully created "+env.Name))
	_ = sendWSClosureMsg(conn, websocket.CloseNormalClosure)
}

func (s *server) handlerDestroyEnv(w http.ResponseWriter, r *http.Request) {
	envName := mux.Vars(r)["name"]

	conn, _ := s.wsUpgrader.Upgrade(w, r, nil)
	defer conn.Close()

	env, err := s.state.Environment(r.Context(), envName)
	if err != nil {
		s.log.Errorf("Failed to fetch env: %v", err)
		_ = sendWSClosureMsg(conn, websocket.CloseInternalServerErr)
		return
	}
	if env == nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Environment does not exist"))
		_ = sendWSClosureMsg(conn, websocket.ClosePolicyViolation)
		return
	}
	if env.Status != environment.StatusProvisioned {
		conn.WriteMessage(websocket.TextMessage, []byte("Environment is not in provisioned state"))
		_ = sendWSClosureMsg(conn, websocket.ClosePolicyViolation)
		return
	}

	// don't destroy the environment if even a single app exists in it
	hasApps, err := s.state.CheckEnvContainsApps(r.Context(), env.Name)
	if err != nil {
		s.log.Errorf("Failed to check if env contains apps: %v", err)
		_ = sendWSClosureMsg(conn, websocket.CloseInternalServerErr)
		return
	}
	if hasApps {
		conn.WriteMessage(
			websocket.TextMessage,
			[]byte("Environment cannot be destroyed because it contains applications"),
		)
		_ = sendWSClosureMsg(conn, websocket.ClosePolicyViolation)
		return
	}

	s.log.WithField("name", envName).Info("Destroying environment")

	env.Status = environment.StatusDestroying
	if err := s.state.UpdateEnvironment(r.Context(), env); err != nil {
		s.log.Errorf("Failed to update env status: %v", err)
		_ = sendWSClosureMsg(conn, websocket.CloseInternalServerErr)
		return
	}
	env.Infra = s.infra

	eventsCh := make(chan environment.Event)
	go env.Destroy(r.Context(), eventsCh)

	for e := range eventsCh {
		if e.Err != nil {
			s.log.Errorf("Failed to destroy environment: %v", e.Err)
			conn.WriteMessage(websocket.TextMessage, []byte("Failed to destroy environment"))
			_ = sendWSClosureMsg(conn, websocket.CloseInternalServerErr)
			return
		}
		s.log.Info(e.Msg)
		conn.WriteMessage(websocket.TextMessage, []byte(e.Msg))
	}

	if err := os.Remove(s.envTfFile(env.Name)); err != nil {
		s.log.Errorf("Failed to delete env TF config file from disk: %v", err)
		_ = sendWSClosureMsg(conn, websocket.CloseInternalServerErr)
		return
	}
	if err := s.state.DeleteEnvironment(r.Context(), envName); err != nil {
		s.log.Errorf("Failed to delete env from state: %v", err)
		_ = sendWSClosureMsg(conn, websocket.CloseInternalServerErr)
		return
	}

	conn.WriteMessage(websocket.TextMessage, []byte("Environment destroyed successfully"))
	_ = sendWSClosureMsg(conn, websocket.CloseNormalClosure)
}

func (s *server) envTfFile(name string) string {
	return path.Join(s.config.DataDir, TerraformDir, name+".tf")
}
