package server

import (
	"github.com/cloudfauj/cloudfauj/infrastructure"
	"github.com/cloudfauj/cloudfauj/state"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path"
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

	ar := r.PathPrefix("/app").Subrouter()
	ar.HandleFunc("/{name}", s.handlerDestroyApp).Methods(http.MethodDelete)
	ar.HandleFunc("/deploy", s.handlerDeployApp)

	dr := r.PathPrefix("/deployment").Subrouter()
	dr.HandleFunc("/{id}", s.handlerGetDeployment).Methods(http.MethodGet)
	dr.HandleFunc("/{id}/logs", s.handlerGetDeploymentLogs).Methods(http.MethodGet)

	// TODO: refactor & take out the /{name} path since its common
	//  Also check if we can write a single controller to handle both plan & apply
	//  since they're same except for TF action.
	er := r.PathPrefix("/environment").Subrouter()
	er.HandleFunc("/create", s.handlerCreateEnv)
	er.HandleFunc("/{name}/destroy", s.handlerDestroyEnv)
	er.HandleFunc("/{name}/plan", s.handlerTFPlanEnv)
	er.HandleFunc("/{name}/apply", s.handlerTFApplyEnv)

	dmr := r.PathPrefix("/domain").Subrouter()
	dmr.HandleFunc("/add", s.handlerAddDomain)
	dmr.HandleFunc("/{name}/delete", s.handlerDeleteDomain)
	dmr.HandleFunc("/{name}/plan", s.handlerTFPlanDomain)
	dmr.HandleFunc("/{name}/apply", s.handlerTFApplyDomain)
}

func (s *server) handlerGetHealthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// writeFiles writes a set of files inside a target directory.
// It takes files as input - a map with filenames as keys and
// their respective contents as values.
func (s *server) writeFiles(dir string, files map[string]string) error {
	for fName, data := range files {
		filepath := path.Join(dir, fName)
		if err := os.WriteFile(filepath, []byte(data), 0666); err != nil {
			return err
		}
	}
	return nil
}
