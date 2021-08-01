package state

import (
	"context"
	"database/sql"
	"github.com/cloudfauj/cloudfauj/deployment"
)

const sqlCreateDeploymentTable = `CREATE TABLE IF NOT EXISTS deployments (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	app VARCHAR(100) NOT NULL,
	env VARCHAR(100) NOT NULL,
	status VARCHAR(40) NOT NULL
)`

func (s *state) Deployment(ctx context.Context, id string) (*deployment.Deployment, error) {
	var d deployment.Deployment

	err := s.db.QueryRowContext(
		ctx, "SELECT * FROM deployments WHERE id = ?", id,
	).Scan(
		&d.Id, &d.App, &d.Environment, &d.Status,
	)
	if err != nil {
		// return nil response without any error if no such env found
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &d, nil
}

// ListDeployments returns a list of all deployments having the specified status
func (s *state) ListDeployments(ctx context.Context, status string) ([]*deployment.Deployment, error) {
	var res []*deployment.Deployment

	rows, err := s.db.QueryContext(ctx, "SELECT * FROM deployments WHERE status = ?", status)
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var d deployment.Deployment
		if err := rows.Scan(&d.Id, &d.App, &d.Environment, &d.Status); err != nil {
			return res, err
		}
		res = append(res, &d)
	}
	err = rows.Err()
	return res, err
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
	q := "UPDATE deployments SET status = ? WHERE id = ?"
	stmt, err := s.db.PrepareContext(ctx, q)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, status, id)
	return err
}
