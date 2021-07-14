package state

import (
	"context"
	"github.com/cloudfauj/cloudfauj/environment"
)

func (s *state) CheckEnvExists(ctx context.Context, name string) (bool, error) {
	return false, nil
}

func (s *state) CreateEnvironment(ctx context.Context, e *environment.Environment) error {
	return nil
}

func (s *state) UpdateEnvironment(ctx context.Context, e *environment.Environment) error {
	return nil
}

func (s *state) ListEnvironments(ctx context.Context) ([]string, error) {
	return []string{}, nil
}

func (s *state) Environment(ctx context.Context, name string) (*environment.Environment, error) {
	// return nil if env doesn't exist
	return nil, nil
}

func (s *state) DeleteEnvironment(context.Context, string) error {
	return nil
}
