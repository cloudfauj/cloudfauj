package state

import (
	"context"
	"github.com/sirupsen/logrus"
)

type State interface {
	CreateEnvironment(ctx context.Context) error
}

type state struct {
	log *logrus.Logger
}

func New(l *logrus.Logger) State {
	return &state{log: l}
}
