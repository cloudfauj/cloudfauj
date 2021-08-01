package state

import (
	"context"
	"github.com/cloudfauj/cloudfauj/deployment"
)

const sqlCreateDeploymentTable = `CREATE TABLE IF NOT EXISTS deployments (
	id INTEGER AUTOINCREMENT,
	app VARCHAR(100) NOT NULL,
	env VARCHAR(100) NOT NULL,
	status VARCHAR(40) NOT NULL
)`

func (s *state) Deployment(ctx context.Context, id string) (*deployment.Deployment, error) {
	// return nil if deployment doesn't exist
	return nil, nil
}

func (s *state) ListDeployments(ctx context.Context, statuses ...string) ([]*deployment.Deployment, error) {
	return nil, nil
}

// CreateDeployment creates a new deployment in state and returns its unique ID
func (s *state) CreateDeployment(ctx context.Context, dep *deployment.Deployment) (int64, error) {
	q := "INSERT INTO deployments(app, env, status) VALUES(?, ?, ?)"
	stmt, err := s.db.PrepareContext(ctx, q)
	if err != nil {
		return 0, err
	}
	res, err := stmt.ExecContext(ctx, dep.App, dep.Environment, dep.Status)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *state) UpdateDeploymentStatus(ctx context.Context, id, status string) error {
	return nil
}
