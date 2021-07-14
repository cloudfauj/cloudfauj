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

func (s *server) handlerDeployApp(w http.ResponseWriter, r *http.Request) {}
