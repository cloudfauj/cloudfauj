package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (s *server) handlerListEnvironments(w http.ResponseWriter, r *http.Request) {}

func (s *server) handlerCreateEnv(w http.ResponseWriter, r *http.Request) {
	s.log.Info("Create env")
	s.log.Info(mux.Vars(r)["name"])
}

func (s *server) handlerDestroyEnv(w http.ResponseWriter, r *http.Request) {}
