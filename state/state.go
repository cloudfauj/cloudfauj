package state

import (
	"context"
	"database/sql"
	"github.com/cloudfauj/cloudfauj/application"
	"github.com/cloudfauj/cloudfauj/deployment"
	"github.com/cloudfauj/cloudfauj/domain"
	"github.com/cloudfauj/cloudfauj/environment"
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
	// CheckEnvContainsApps returns true if the given environment contains even a single application
	CheckEnvContainsApps(context.Context, string) (bool, error)

	Deployment(context.Context, string) (*deployment.Deployment, error)
	ListDeployments(context.Context, string) ([]*deployment.Deployment, error)
	CreateDeployment(context.Context, *deployment.Deployment) (int64, error)
	UpdateDeploymentStatus(context.Context, string, string) error

	//CheckAppExists(context.Context, string, string) (bool, error)
	CreateApp(context.Context, *application.Application, string) error
	UpdateApp(context.Context, *application.Application, string) error
	App(context.Context, string, string) (*application.Application, error)
	DeleteApp(context.Context, string, string) error

	AddDomain(context.Context, *domain.Domain) error
	CheckDomainExists(context.Context, string) (bool, error)
	DeleteDomain(context.Context, string) error
	ListDomains(context.Context) ([]string, error)
}

type state struct {
	log *logrus.Logger
	db  *sql.DB
}

func New(l *logrus.Logger, db *sql.DB) State {
	return &state{log: l, db: db}
}
