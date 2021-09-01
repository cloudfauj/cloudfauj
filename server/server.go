package server

import (
	"github.com/cloudfauj/cloudfauj/infrastructure"
	"github.com/cloudfauj/cloudfauj/state"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
)

const ApiV1Prefix = "/v1"

// Event represent a server event
type Event struct {
	Msg string
	Err error
}

type server struct {
	config     *Config
	log        *logrus.Logger
	infra      *infrastructure.Infrastructure
	state      state.State
	wsUpgrader *websocket.Upgrader
	*mux.Router
}

func New(c *Config, l *logrus.Logger, s state.State, i *infrastructure.Infrastructure) http.Handler {
	srv := &server{
		config:     c,
		log:        l,
		infra:      i,
		state:      s,
		wsUpgrader: &websocket.Upgrader{},
		Router:     mux.NewRouter(),
	}
	setupV1Routes(srv)
	return srv
}

func setupV1Routes(s *server) {
	r := s.Router.PathPrefix(ApiV1Prefix).Subrouter()

	r.HandleFunc("/health", s.handlerGetHealthcheck).Methods(http.MethodGet)
	r.HandleFunc("/environments", s.handlerListEnvironments).Methods(http.MethodGet)
	r.HandleFunc("/deployments", s.handlerListDeployments).Methods(http.MethodGet)
	r.HandleFunc("/domains", s.handlerListDomains).Methods(http.MethodGet)

	appR := r.PathPrefix("/app").Subrouter()
	appR.HandleFunc("/{name}", s.handlerDestroyApp).Methods(http.MethodDelete)
	appR.HandleFunc("/deploy", s.handlerDeployApp)

	depR := r.PathPrefix("/deployment").Subrouter()
	depR.HandleFunc("/{id}", s.handlerGetDeployment).Methods(http.MethodGet)
	depR.HandleFunc("/{id}/logs", s.handlerGetDeploymentLogs).Methods(http.MethodGet)

	envR := r.PathPrefix("/environment").Subrouter()
	envR.HandleFunc("/create", s.handlerCreateEnv)
	envR.HandleFunc("/{name}/destroy", s.handlerDestroyEnv)

	domainR := r.PathPrefix("/domain").Subrouter()
	domainR.HandleFunc("/{name}/add", s.handlerAddDomain)
	domainR.HandleFunc("/{name}/delete", s.handlerDeleteDomain)
}

func (s *server) handlerGetHealthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
