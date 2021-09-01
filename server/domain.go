package server

import (
	"encoding/json"
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
	nsRecords, err := s.infra.CreateDomain(r.Context(), name, tf, f)
	if err != nil {
		s.log.Errorf("Failed to provision domain infrastructure: %v", err)
		conn.sendFailureISE()
		return
	}

	conn.sendTextMsg("NS Records to be configured for " + name)
	for _, r := range nsRecords {
		conn.sendTextMsg(r)
	}
	conn.sendSuccess("Domain infrastructure created successfully")
}

func (s *server) handlerDeleteDomain(w http.ResponseWriter, r *http.Request) {
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
	if !exists {
		conn.sendSuccess("Domain doesn't exist, nothing to do")
		return
	}

	// TODO: Before proceeding with destruction, do we want to check if
	//  the infra is being relied on? eg- is the ACM cert being used by
	//  any load balancer?

	conn.sendTextMsg("Destroying Terraform infrastructure")

	tf, err := s.infra.NewTerraform(s.domainTFDir(name))
	if err != nil {
		s.log.Error(err)
		conn.sendFailureISE()
		return
	}
	if err := s.infra.DeleteDomain(r.Context(), tf); err != nil {
		s.log.Errorf("Failed to destroy domain infrastructure: %v", err)
		conn.sendFailureISE()
		return
	}
	if err := os.RemoveAll(s.domainTFDir(name)); err != nil {
		s.log.Errorf("Failed to delete domain TF config from disk: %v", err)
		conn.sendFailureISE()
		return
	}

	conn.sendTextMsg(fmt.Sprintf("De-registering %s from state", name))
	if err := s.state.DeleteDomain(r.Context(), name); err != nil {
		s.log.WithField("name", name).Errorf("Failed to delete domain from state: %v", err)
		conn.sendFailureISE()
		return
	}

	conn.sendSuccess("Domain deleted successfully")
}

func (s *server) handlerListDomains(w http.ResponseWriter, r *http.Request) {
	res, err := s.state.ListDomains(r.Context())
	if err != nil {
		s.log.Errorf("Failed to list domains from state: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsonRes, _ := json.Marshal(res)
	_, _ = w.Write(jsonRes)
}

func (s *server) domainTFDir(name string) string {
	return path.Join(s.config.TerraformDomainsDir(), name)
}

func (s *server) domainTFFile(name string) string {
	return path.Join(s.domainTFDir(name), s.config.terraformConfigFile)
}
