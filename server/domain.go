package server

import (
	"encoding/json"
	"fmt"
	"github.com/cloudfauj/cloudfauj/domain"
	"github.com/cloudfauj/cloudfauj/wsmanager"
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
	conn := &wsmanager.WSManager{Conn: wsConn}

	var d *domain.Domain
	if err := conn.ReadJSON(&d); err != nil {
		s.log.Errorf("Failed to read domain config: %v", err)
		conn.SendFailureISE()
		return
	}
	if err := d.CheckIsValid(); err != nil {
		conn.SendFailure(
			fmt.Sprintf("Invalid domain configuration: %v", err),
			websocket.CloseInvalidFramePayloadData,
		)
		return
	}

	exists, err := s.state.CheckDomainExists(r.Context(), d.Name)
	if err != nil {
		s.log.WithField("name", d.Name).Errorf("Failed to check if domain exists: %v", err)
		conn.SendFailureISE()
		return
	}
	if exists {
		conn.SendFailure("Domain already exists", websocket.ClosePolicyViolation)
		return
	}

	conn.SendTextMsg(fmt.Sprintf("Registering %s in state", d.Name))
	if err := s.state.AddDomain(r.Context(), d); err != nil {
		s.log.WithField("name", d.Name).Errorf("Failed to add domain to state: %v", err)
		conn.SendFailureISE()
		return
	}

	conn.SendTextMsg("Generating Terraform configuration")

	// Get the terraform config filenames and their contents
	tfConfigs, err := s.infra.DomainTFConfig(d)
	if err != nil {
		s.log.Errorf("Failed to generate terraform configurations for domain: %v", err)
		conn.SendFailureISE()
		return
	}

	// Create the domain terraform module on disk
	dir := s.domainTFDir(d.Name)
	if err := os.Mkdir(dir, 0755); err != nil {
		s.log.Errorf("Failed to create directory for domain: %v", err)
		conn.SendFailureISE()
		return
	}
	if err := s.writeFiles(dir, tfConfigs); err != nil {
		s.log.Errorf("Failed to write terraform configs for domain: %v", err)
		conn.SendFailureISE()
		return
	}

	conn.SendTextMsg("Provisioning infrastructure")

	// Provision domain infrastructure by invoking terraform
	tf, err := s.infra.NewTerraform(dir, conn)
	if err != nil {
		s.log.Error(err)
		conn.SendFailureISE()
		return
	}

	nsRecords, err := s.infra.CreateDomain(r.Context(), tf)
	if err != nil {
		s.log.Errorf("Failed to provision domain infrastructure: %v", err)
		conn.SendFailureISE()
		return
	}
	conn.SendTextMsg("NS Records to be configured for " + d.Name)
	for _, r := range nsRecords {
		conn.SendTextMsg(r)
	}

	conn.SendSuccess("Domain infrastructure created successfully")
}

func (s *server) handlerDeleteDomain(w http.ResponseWriter, r *http.Request) {
	wsConn, err := s.wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log.Errorf("Failed to upgrade websocket connection: %v", err)
		return
	}
	defer wsConn.Close()
	conn := &wsmanager.WSManager{Conn: wsConn}

	name := mux.Vars(r)["name"]
	exists, err := s.state.CheckDomainExists(r.Context(), name)
	if err != nil {
		s.log.WithField("name", name).Errorf("Failed to check if domain exists: %v", err)
		conn.SendFailureISE()
		return
	}
	if !exists {
		conn.SendSuccess("Domain doesn't exist, nothing to do")
		return
	}

	// TODO: Abort if domain being used by any environments
	conn.SendTextMsg("Destroying infrastructure")

	dir := s.domainTFDir(name)
	tf, err := s.infra.NewTerraform(dir, conn)
	if err != nil {
		s.log.Error(err)
		conn.SendFailureISE()
		return
	}
	if err := s.infra.DeleteDomain(r.Context(), tf); err != nil {
		s.log.Errorf("Failed to destroy domain infrastructure: %v", err)
		conn.SendFailureISE()
		return
	}
	if err := os.RemoveAll(dir); err != nil {
		s.log.Errorf("Failed to delete domain TF config from disk: %v", err)
		conn.SendFailureISE()
		return
	}

	conn.SendTextMsg(fmt.Sprintf("De-registering %s from state", name))
	if err := s.state.DeleteDomain(r.Context(), name); err != nil {
		s.log.WithField("name", name).Errorf("Failed to delete domain from state: %v", err)
		conn.SendFailureISE()
		return
	}

	conn.SendSuccess("Domain deleted successfully")
}

func (s *server) handlerTFPlanDomain(w http.ResponseWriter, r *http.Request) {
	wsConn, err := s.wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log.Errorf("Failed to upgrade websocket connection: %v", err)
		return
	}
	defer wsConn.Close()
	conn := &wsmanager.WSManager{Conn: wsConn}

	name := mux.Vars(r)["name"]
	exists, err := s.state.CheckDomainExists(r.Context(), name)
	if err != nil {
		s.log.WithField("name", name).Errorf("Failed to check if domain exists: %v", err)
		conn.SendFailureISE()
		return
	}
	if !exists {
		conn.SendSuccess("Domain doesn't exist, nothing to do")
		return
	}

	conn.SendTextMsg(fmt.Sprintf("Running Terraform Plan over %s infrastructure configuration", name))

	dir := s.domainTFDir(name)
	tf, err := s.infra.NewTerraform(dir, conn)
	if err != nil {
		s.log.Error(err)
		conn.SendFailureISE()
		return
	}

	// TODO: Recursively plan all other infra modules that are dependent on
	//  this domain module (eg- environments using this domain).
	if _, err := s.infra.PlanDomain(r.Context(), tf); err != nil {
		s.log.Errorf("Failed to plan domain infrastructure: %v", err)
		conn.SendFailureISE()
		return
	}
	conn.SendSuccess("Done")
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

func (s *server) domainTFStateFile(name string) string {
	return path.Join(s.domainTFDir(name), s.config.terraformStateFile)
}
