package server

import (
	"github.com/cloudfauj/cloudfauj/state"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
)

// config fields should be public
type Config struct {
	// DataDir is the base directory inside which Cloudfauj server
	// stores all its data.
	// To restore Cloudfauj server on to a new server, restoring a
	// backup of this dir and running the server is enough.
	DataDir string `mapstructure:"data_dir"`
}

type server struct {
	config     *Config
	log        *logrus.Logger
	state      state.State
	wsUpgrader *websocket.Upgrader
	*mux.Router
}

const ApiV1Prefix = "/v1"

const (
	DeploymentsDir     = "deployments"
	ApplicationsDir    = "applications"
	ApplicationsEnvDir = "env"
	LogFileBasename    = "logs.txt"
)

func New(c *Config, l *logrus.Logger, s state.State) http.Handler {
	srv := &server{
		config:     c,
		log:        l,
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
