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

// State manages all structured data persisted on disk for Cloudfauj Server
type State interface {
	// Migrate runs DB migrations to prepare all tables for the server to work with
	Migrate(context.Context) error

	CheckEnvExists(context.Context, string) (bool, error)
	CreateEnvironment(context.Context, *environment.Environment) error
	UpdateEnvironment(context.Context, *environment.Environment) error
	ListEnvironments(context.Context) ([]string, error)
	Environment(context.Context, string) (*environment.Environment, error)
	DeleteEnvironment(context.Context, string) error

	Deployment(context.Context, string) (*deployment.Deployment, error)
	ListDeployments(context.Context, string) ([]*deployment.Deployment, error)
	CreateDeployment(context.Context, *deployment.Deployment) (int64, error)
	UpdateDeploymentStatus(context.Context, string, string) error

	CheckAppExists(context.Context, string) (bool, error)
	CreateApp(context.Context, *application.Application) error
	UpdateApp(context.Context, *application.Application) error
	App(context.Context, string) (*application.Application, error)
	DeleteApp(context.Context, string) error

	CreateAppInfra(context.Context, *infrastructure.AppInfra) error
	UpdateAppInfra(context.Context, *infrastructure.AppInfra) error
	AppInfra(context.Context, string) (*infrastructure.AppInfra, error)
	DeleteAppInfra(context.Context, string) error
}

type state struct {
	log *logrus.Logger
	db  *sql.DB
}

func New(l *logrus.Logger, db *sql.DB) State {
	return &state{log: l, db: db}
}
