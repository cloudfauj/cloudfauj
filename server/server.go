package server

import (
	"github.com/cloudfauj/cloudfauj/state"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

// todo: evaluate whether we even need this
// config fields should be public
type Config struct{}

type server struct {
	config *Config
	log    *logrus.Logger
	state  state.State
	*mux.Router
}

const ApiV1Prefix = "/v1"

func New(c *Config, l *logrus.Logger, s state.State) http.Handler {
	srv := &server{
		config: c,
		log:    l,
		state:  s,
		Router: mux.NewRouter(),
	}
	setupV1Routes(srv)
	return srv
}

func setupV1Routes(s *server) {
	r := s.Router.PathPrefix(ApiV1Prefix).Subrouter()
	r.HandleFunc("/health", s.handlerGetHealthcheck).Methods(http.MethodGet)
	r.HandleFunc("/environments", s.handlerListEnvironments).Methods(http.MethodGet)
	r.HandleFunc("/deployments", s.handlerListDeployments).Methods(http.MethodGet)

	appR := r.PathPrefix("/app").Subrouter()
	appR.HandleFunc("/{name}/logs", s.handlerGetAppLogs).Methods(http.MethodGet)
	appR.HandleFunc("/deploy", s.handlerDeployApp)

	depR := r.PathPrefix("/deployment").Subrouter()
	depR.HandleFunc("/{id}", s.handlerGetDeployment).Methods(http.MethodGet)
	depR.HandleFunc("/{id}/logs", s.handlerGetDeploymentLogs).Methods(http.MethodGet)

	envR := r.PathPrefix("/environment").Subrouter()
	envR.HandleFunc("/{name}/create", s.handlerCreateEnv)
	envR.HandleFunc("/{name}/destroy", s.handlerDestroyEnv)
}