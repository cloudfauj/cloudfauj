package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func (s *server) handlerGetDeployment(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	d, err := s.state.GetDeployment(r.Context(), id)
	if err != nil {
		s.log.Errorf("Failed to fetch deployment from state: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if d == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsonRes, _ := json.Marshal(d)
	_, _ = w.Write(jsonRes)
}

func (s *server) handlerGetDeploymentLogs(w http.ResponseWriter, r *http.Request) {}

func (s *server) handlerListDeployments(w http.ResponseWriter, r *http.Request) {}
