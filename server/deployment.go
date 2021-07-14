package server

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io/fs"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
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

func (s *server) handlerGetDeploymentLogs(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	f := path.Join(s.config.DataDir, DeploymentsDir, id, LogFileBasename)

	s.log.WithField("path", f).Info("Fetching deployment logs")

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

func (s *server) handlerListDeployments(w http.ResponseWriter, r *http.Request) {}
