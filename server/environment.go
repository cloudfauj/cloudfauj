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
	wsConn, err := s.wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log.Errorf("Failed to upgrade websocket connection: %v", err)
		return
	}
	defer wsConn.Close()
	conn := wsManager{wsConn}

	var env *environment.Environment
	if err := conn.ReadJSON(&env); err != nil {
		s.log.Errorf("Failed to read environment config: %v", err)
		conn.sendFailureISE()
		return
	}
	if err := env.CheckIsValid(); err != nil {
		conn.sendFailure(
			fmt.Sprintf("Invalid environment config: %v", err),
			websocket.CloseInvalidFramePayloadData,
		)
		return
	}

	ok, err := s.state.CheckEnvExists(r.Context(), env.Name)
	if err != nil {
		s.log.Errorf("Failed to check if env exists: %v", err)
		conn.sendFailureISE()
		return
	}
	if ok {
		conn.sendFailure("Environment already exists", websocket.ClosePolicyViolation)
		return
	}

	if env.DomainEnabled() {
		exists, err := s.state.CheckDomainExists(r.Context(), env.Domain)
		if err != nil {
			s.log.Errorf("Failed to check if domain to use for env exists: %v", err)
			conn.sendFailureISE()
			return
		}
		if !exists {
			conn.sendFailure("Specified domain does not exist in the system", websocket.ClosePolicyViolation)
			return
		}
	}

	s.log.WithField("name", env.Name).Info("Creating new environment")

	env.Status = environment.StatusProvisioning
	if err := s.state.CreateEnvironment(r.Context(), env); err != nil {
		s.log.Errorf("Failed to store env info in state: %v", err)
		conn.sendFailureISE()
		return
	}
	conn.sendTextMsg("Registered in state")

	if err := os.Mkdir(s.envTfDir(env.Name), 0755); err != nil {
		s.log.Errorf("Failed to create directory for env: %v", err)
		conn.sendFailureISE()
		return
	}
	f, err := os.OpenFile(s.envTfFile(env.Name), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		s.log.Errorf("Failed to create Terraform config file for env: %v", err)
		conn.sendFailureISE()
		return
	}
	defer f.Close()

	tf, err := s.infra.NewTerraform(s.envTfDir(env.Name))
	if err != nil {
		s.log.Error(err)
		conn.sendFailureISE()
		return
	}

	conn.sendTextMsg("Applying Terraform configuration")
	// TODO: Is there a more idiomatic way than to supply domain tf state file?
	err = s.infra.CreateEnvironment(r.Context(), env, s.domainTFStateFile(env.Domain), tf, f)
	if err != nil {
		s.log.Errorf("Failed to provision environment: %v", err)
		conn.sendFailureISE()
		return
	}

	env.Status = environment.StatusProvisioned
	if err := s.state.UpdateEnvStatus(r.Context(), env.Name, env.Status); err != nil {
		s.log.Errorf("Failed to update env info in state: %v", err)
		conn.sendFailureISE()
		return
	}
	conn.sendSuccess("Successfully created " + env.Name)
}

func (s *server) handlerDestroyEnv(w http.ResponseWriter, r *http.Request) {
	wsConn, err := s.wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log.Errorf("Failed to upgrade websocket connection: %v", err)
		return
	}
	defer wsConn.Close()
	conn := wsManager{wsConn}

	envName := mux.Vars(r)["name"]

	env, err := s.state.Environment(r.Context(), envName)
	if err != nil {
		s.log.Errorf("Failed to fetch env: %v", err)
		conn.sendFailureISE()
		return
	}
	if env == nil {
		conn.sendFailure("Environment does not exist", websocket.ClosePolicyViolation)
		return
	}
	if env.Status != environment.StatusProvisioned {
		conn.sendFailure("Environment is not in provisioned state", websocket.ClosePolicyViolation)
		return
	}

	// don't destroy the environment if even a single app exists in it
	hasApps, err := s.state.CheckEnvContainsApps(r.Context(), env.Name)
	if err != nil {
		s.log.Errorf("Failed to check if env contains apps: %v", err)
		conn.sendFailureISE()
		return
	}
	if hasApps {
		conn.sendFailure(
			"Environment cannot be destroyed because it contains applications",
			websocket.ClosePolicyViolation,
		)
		return
	}

	s.log.WithField("name", envName).Info("Destroying environment")

	env.Status = environment.StatusDestroying
	if err := s.state.UpdateEnvStatus(r.Context(), env.Name, env.Status); err != nil {
		s.log.Errorf("Failed to update env status: %v", err)
		conn.sendFailureISE()
		return
	}

	tf, err := s.infra.NewTerraform(s.envTfDir(env.Name))
	if err != nil {
		s.log.Error(err)
		conn.sendFailureISE()
		return
	}

	conn.sendTextMsg("Destroying Terraform infrastructure")
	if err := s.infra.DestroyEnvironment(r.Context(), tf); err != nil {
		s.log.Errorf("Failed to destroy environment: %v", err)
		conn.sendFailureISE()
		return
	}
	if err := os.RemoveAll(s.envTfDir(env.Name)); err != nil {
		s.log.Errorf("Failed to delete env TF config file from disk: %v", err)
		conn.sendFailureISE()
		return
	}
	if err := s.state.DeleteEnvironment(r.Context(), envName); err != nil {
		s.log.Errorf("Failed to delete env from state: %v", err)
		conn.sendFailureISE()
		return
	}

	conn.sendSuccess("Environment destroyed successfully")
}

func (s *server) envTfDir(name string) string {
	return path.Join(s.config.TerraformDir(), name)
}

func (s *server) envTfFile(name string) string {
	return path.Join(s.envTfDir(name), s.config.terraformConfigFile)
}

func (s *server) envTfStateFile(name string) string {
	return path.Join(s.envTfDir(name), s.config.terraformStateFile)
}
