package server

import (
	"github.com/cloudfauj/cloudfauj/infrastructure"
	"github.com/cloudfauj/cloudfauj/state"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
)

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

const ApiV1Prefix = "/v1"

func New(c *Config, l *logrus.Logger, s state.State) http.Handler {
	srv := &server{
		config:     c,
		log:        l,
		infra:      infrastructure.New(),
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

	appR := r.PathPrefix("/app").Subrouter()
	appR.HandleFunc("/{name}/logs", s.handlerGetAppLogs).Methods(http.MethodGet)
	appR.HandleFunc("/deploy", s.handlerDeployApp)

	depR := r.PathPrefix("/deployment").Subrouter()
	depR.HandleFunc("/{id}", s.handlerGetDeployment).Methods(http.MethodGet)
	depR.HandleFunc("/{id}/logs", s.handlerGetDeploymentLogs).Methods(http.MethodGet)

	envR := r.PathPrefix("/environment").Subrouter()
	envR.HandleFunc("/create", s.handlerCreateEnv)
	envR.HandleFunc("/{name}/destroy", s.handlerDestroyEnv)
}

func sendWSClosureMsg(conn *websocket.Conn, code int) error {
	return conn.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(code, ""),
	)
}

func (s *server) handlerGetHealthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
