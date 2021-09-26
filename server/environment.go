package server

import (
	"encoding/json"
	"fmt"
	"github.com/cloudfauj/cloudfauj/environment"
	"github.com/cloudfauj/cloudfauj/wsmanager"
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
		s.log.Errorf("Failed to upgrade wsmanager connection: %v", err)
		return
	}
	defer wsConn.Close()
	conn := &wsmanager.WSManager{Conn: wsConn}

	var env *environment.Environment
	if err := conn.ReadJSON(&env); err != nil {
		s.log.Errorf("Failed to read environment config: %v", err)
		conn.SendFailureISE()
		return
	}
	if err := env.CheckIsValid(); err != nil {
		conn.SendFailure(
			fmt.Sprintf("Invalid environment config: %v", err),
			websocket.CloseInvalidFramePayloadData,
		)
		return
	}

	ok, err := s.state.CheckEnvExists(r.Context(), env.Name)
	if err != nil {
		s.log.Errorf("Failed to check if env exists: %v", err)
		conn.SendFailureISE()
		return
	}
	if ok {
		conn.SendFailure("Environment already exists", websocket.ClosePolicyViolation)
		return
	}

	if env.DomainEnabled() {
		exists, err := s.state.CheckDomainExists(r.Context(), env.Domain)
		if err != nil {
			s.log.Errorf("Failed to check if domain to use for env exists: %v", err)
			conn.SendFailureISE()
			return
		}
		if !exists {
			conn.SendFailure("Specified domain does not exist in the system", websocket.ClosePolicyViolation)
			return
		}
	}

	s.log.WithField("name", env.Name).Info("Creating new environment")

	env.Status = environment.StatusProvisioning
	if err := s.state.CreateEnvironment(r.Context(), env); err != nil {
		s.log.Errorf("Failed to store env info in state: %v", err)
		conn.SendFailureISE()
		return
	}
	conn.SendTextMsg("Registered in state")
	conn.SendTextMsg("Generating Terraform configuration")

	tfConfigs, err := s.infra.EnvTFConfig(r.Context(), env, s.domainTFStateFile(env.Domain))
	if err != nil {
		s.log.Errorf("Failed to generate terraform configurations for env: %v", err)
		conn.SendFailureISE()
		return
	}
	dir := s.envTfDir(env.Name)
	if err := os.Mkdir(dir, 0755); err != nil {
		s.log.Errorf("Failed to create directory for env: %v", err)
		conn.SendFailureISE()
		return
	}
	if err := s.writeFiles(dir, tfConfigs); err != nil {
		s.log.Errorf("Failed to write terraform configs for environment: %v", err)
		conn.SendFailureISE()
		return
	}

	conn.SendTextMsg("Provisioning infrastructure")

	tf, err := s.infra.NewTerraform(dir)
	if err != nil {
		s.log.Error(err)
		conn.SendFailureISE()
		return
	}
	err = s.infra.CreateEnvironment(r.Context(), tf)
	if err != nil {
		s.log.Errorf("Failed to provision environment: %v", err)
		conn.SendFailureISE()
		return
	}

	env.Status = environment.StatusProvisioned
	if err := s.state.UpdateEnvStatus(r.Context(), env.Name, env.Status); err != nil {
		s.log.Errorf("Failed to update env info in state: %v", err)
		conn.SendFailureISE()
		return
	}
	conn.SendSuccess("Successfully created " + env.Name)
}

func (s *server) handlerDestroyEnv(w http.ResponseWriter, r *http.Request) {
	wsConn, err := s.wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log.Errorf("Failed to upgrade wsmanager connection: %v", err)
		return
	}
	defer wsConn.Close()
	conn := &wsmanager.WSManager{Conn: wsConn}

	envName := mux.Vars(r)["name"]

	env, err := s.state.Environment(r.Context(), envName)
	if err != nil {
		s.log.Errorf("Failed to fetch env: %v", err)
		conn.SendFailureISE()
		return
	}
	if env == nil {
		conn.SendFailure("Environment does not exist", websocket.ClosePolicyViolation)
		return
	}
	if env.Status != environment.StatusProvisioned {
		conn.SendFailure("Environment is not in provisioned state", websocket.ClosePolicyViolation)
		return
	}

	// don't destroy the environment if even a single app exists in it
	hasApps, err := s.state.CheckEnvContainsApps(r.Context(), env.Name)
	if err != nil {
		s.log.Errorf("Failed to check if env contains apps: %v", err)
		conn.SendFailureISE()
		return
	}
	if hasApps {
		conn.SendFailure(
			"Environment cannot be destroyed because it contains applications",
			websocket.ClosePolicyViolation,
		)
		return
	}

	s.log.WithField("name", envName).Info("Destroying environment")

	env.Status = environment.StatusDestroying
	if err := s.state.UpdateEnvStatus(r.Context(), env.Name, env.Status); err != nil {
		s.log.Errorf("Failed to update env status: %v", err)
		conn.SendFailureISE()
		return
	}

	tf, err := s.infra.NewTerraform(s.envTfDir(env.Name))
	if err != nil {
		s.log.Error(err)
		conn.SendFailureISE()
		return
	}

	conn.SendTextMsg("Destroying Terraform infrastructure")
	if err := s.infra.DestroyEnvironment(r.Context(), tf); err != nil {
		s.log.Errorf("Failed to destroy environment: %v", err)
		conn.SendFailureISE()
		return
	}
	if err := os.RemoveAll(s.envTfDir(env.Name)); err != nil {
		s.log.Errorf("Failed to delete env TF config file from disk: %v", err)
		conn.SendFailureISE()
		return
	}
	if err := s.state.DeleteEnvironment(r.Context(), envName); err != nil {
		s.log.Errorf("Failed to delete env from state: %v", err)
		conn.SendFailureISE()
		return
	}

	conn.SendSuccess("Environment destroyed successfully")
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
