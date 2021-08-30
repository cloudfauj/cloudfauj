package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
	"path"
)

func (s *server) handlerAddDomain(w http.ResponseWriter, r *http.Request) {
	wsConn, err := s.wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log.Errorf("Failed to upgrade websocket connection: %v", err)
		return
	}
	defer wsConn.Close()
	conn := wsManager{wsConn}

	name := mux.Vars(r)["name"]
	exists, err := s.state.CheckDomainExists(r.Context(), name)
	if err != nil {
		s.log.WithField("name", name).Errorf("Failed to check if domain exists: %v", err)
		conn.sendFailureISE()
		return
	}
	if exists {
		conn.sendFailure("Domain already exists", websocket.ClosePolicyViolation)
		return
	}

	conn.sendTextMsg(fmt.Sprintf("Registering %s in state", name))
	if err := s.state.AddDomain(r.Context(), name); err != nil {
		s.log.WithField("name", name).Errorf("Failed to add domain to state: %v", err)
		conn.sendFailureISE()
		return
	}

	conn.sendTextMsg("Setting up Terraform configuration")

	if err := os.Mkdir(s.domainTFDir(name), 0755); err != nil {
		s.log.Errorf("Failed to create directory for domain: %v", err)
		conn.sendFailureISE()
		return
	}
	f, err := os.OpenFile(s.domainTFFile(name), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		s.log.Errorf("Failed to create Terraform config file for domain: %v", err)
		conn.sendFailureISE()
		return
	}
	defer f.Close()

	tf, err := s.infra.NewTerraform(s.domainTFDir(name))
	if err != nil {
		s.log.Error(err)
		conn.sendFailureISE()
		return
	}

	conn.sendTextMsg("Applying Terraform configuration")
	nsRecords, err := s.infra.AddDomain(r.Context(), name, tf, f)
	if err != nil {
		s.log.Errorf("Failed to provision domain infrastructure: %v", err)
		conn.sendFailureISE()
		return
	}

	conn.sendTextMsg("Domain added successfully")
	conn.sendTextMsg(fmt.Sprintf("NS Records: %v", nsRecords))
}

func (s *server) handlerDeleteDomain(w http.ResponseWriter, r *http.Request) {}

func (s *server) domainTFDir(name string) string {
	return path.Join(s.config.TerraformDomainsDir(), name)
}

func (s *server) domainTFFile(name string) string {
	return path.Join(s.domainTFDir(name), s.config.terraformConfigFile)
}
