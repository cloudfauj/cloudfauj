package state

import (
	"context"
	"database/sql"
	"github.com/cloudfauj/cloudfauj/environment"
)

func (s *state) CheckEnvExists(ctx context.Context, name string) (bool, error) {
	var res string
	err := s.db.QueryRow("SELECT name FROM environments WHERE name = ?", name).Scan(&res)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
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
