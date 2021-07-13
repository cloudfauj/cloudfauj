package state

import (
	"context"
	"github.com/cloudfauj/cloudfauj/deployment"
	"github.com/cloudfauj/cloudfauj/environment"
	"github.com/sirupsen/logrus"
)

type State interface {
	CheckEnvExists(context.Context, string) (bool, error)
	CreateEnvironment(context.Context, *environment.Environment) error
	UpdateEnvironment(context.Context, *environment.Environment) error
	ListEnvironments(context.Context) ([]string, error)
	GetEnvironment(context.Context, string) (*environment.Environment, error)
	DeleteEnvironment(context.Context, string) error

	GetDeployment(context.Context, string) (*deployment.Deployment, error)
}

type state struct {
	log *logrus.Logger
}

func New(l *logrus.Logger) State {
	return &state{log: l}
}
