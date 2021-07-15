package state

import (
	"context"
	"github.com/cloudfauj/cloudfauj/application"
)

func (s *state) CheckAppExists(ctx context.Context, name string) (bool, error) {
	return false, nil
}

func (s *state) CreateApp(context.Context, *application.Application) error {
	return nil
}

func (s *state) UpdateApp(context.Context, *application.Application) error {
	return nil
}

func (s *state) App(ctx context.Context, name string) (*application.Application, error) {
	// return nil if app doesn't exist in state
	return nil, nil
}
