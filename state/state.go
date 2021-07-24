package state

import (
	"context"
	"database/sql"
	"github.com/cloudfauj/cloudfauj/application"
	"github.com/cloudfauj/cloudfauj/deployment"
	"github.com/cloudfauj/cloudfauj/environment"
	"github.com/cloudfauj/cloudfauj/infrastructure"
	"github.com/sirupsen/logrus"
)

// State manages all structured & blob data persisted on disk
// for Cloudfauj Server.
type State interface {
	CheckEnvExists(context.Context, string) (bool, error)
	CreateEnvironment(context.Context, *environment.Environment) error
	UpdateEnvironment(context.Context, *environment.Environment) error
	ListEnvironments(context.Context) ([]string, error)
	Environment(context.Context, string) (*environment.Environment, error)
	DeleteEnvironment(context.Context, string) error

	Deployment(context.Context, string) (*deployment.Deployment, error)
	ListDeployments(context.Context, ...string) ([]*deployment.Deployment, error)
	CreateDeployment(context.Context, *deployment.Deployment) (string, error)
	UpdateDeploymentStatus(context.Context, string) error

	CheckAppExists(context.Context, string) (bool, error)
	CreateApp(context.Context, *application.Application) error
	UpdateApp(context.Context, *application.Application) error
	App(context.Context, string) (*application.Application, error)

	CreateAppInfra(context.Context, *infrastructure.AppInfra) error
	UpdateAppInfra(context.Context, *infrastructure.AppInfra) error
	AppInfra(context.Context, string) (*infrastructure.AppInfra, error)
}

type state struct {
	log *logrus.Logger
	db  *sql.DB
}

func New(l *logrus.Logger, db *sql.DB) State {
	return &state{log: l, db: db}
}
