package state

import (
	"context"
	"github.com/cloudfauj/cloudfauj/environment"
	"github.com/sirupsen/logrus"
)

type State interface {
	CheckEnvExists(context.Context, string) (bool, error)
	CreateEnvironment(context.Context, *environment.Environment) error
	UpdateEnvironment(context.Context, *environment.Environment) error
	ListEnvironments(context.Context) ([]string, error)
}

type state struct {
	log *logrus.Logger
}

func New(l *logrus.Logger) State {
	return &state{log: l}
}
