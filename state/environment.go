package state

import (
	"context"
	"database/sql"
	"github.com/cloudfauj/cloudfauj/environment"
)

const sqlCreateEnvTable = `CREATE TABLE IF NOT EXISTS environments (
	name VARCHAR(100) NOT NULL PRIMARY KEY,
	status VARCHAR(25) NOT NULL,
	vpc_id VARCHAR(40),
	internet_gateway VARCHAR(50),
	ecs_cluster VARCHAR(100),
	task_exec_iam_role VARCHAR(200),
	compute_subnet VARCHAR(100)
)`

func (s *state) CheckEnvExists(ctx context.Context, name string) (bool, error) {
	var res string
	err := s.db.QueryRowContext(ctx, "SELECT name FROM environments WHERE name = ?", name).Scan(&res)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *state) CreateEnvironment(ctx context.Context, e *environment.Environment) error {
	q := `INSERT INTO environments(
	name,
	status,
	vpc_id,
	internet_gateway,
	ecs_cluster,
	task_exec_iam_role,
	compute_subnet
) VALUES(?, ?, ?, ?, ?, ?, ?)`

	stmt, err := s.db.PrepareContext(ctx, q)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(
		ctx,
		e.Name,
		e.Status,
		e.Res.VpcId,
		e.Res.InternetGateway,
		e.Res.ECSCluster,
		e.Res.TaskExecIAMRole,
		e.Res.ComputeSubnet,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *state) UpdateEnvironment(ctx context.Context, e *environment.Environment) error {
	q := `UPDATE environments
SET
	status = ?,
	vpc_id = ?,
	internet_gateway = ?,
	ecs_cluster = ?,
	task_exec_iam_role = ?,
	compute_subnet = ?
WHERE name = ?`

	stmt, err := s.db.PrepareContext(ctx, q)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(
		ctx,
		e.Status,
		e.Res.VpcId,
		e.Res.InternetGateway,
		e.Res.ECSCluster,
		e.Res.TaskExecIAMRole,
		e.Res.ComputeSubnet,
		e.Name,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *state) ListEnvironments(ctx context.Context) ([]string, error) {
	var res []string

	rows, err := s.db.QueryContext(ctx, "SELECT name FROM environments")
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return res, err
		}
		res = append(res, name)
	}
	if err = rows.Err(); err != nil {
		return res, err
	}

	return res, nil
}

func (s *state) Environment(ctx context.Context, name string) (*environment.Environment, error) {
	e := &environment.Environment{Res: &environment.Resources{}}
	err := s.db.QueryRowContext(
		ctx, "SELECT * FROM environments WHERE name = ?", name,
	).Scan(
		&e.Name,
		&e.Status,
		&e.Res.VpcId,
		&e.Res.InternetGateway,
		&e.Res.ECSCluster,
		&e.Res.TaskExecIAMRole,
		&e.Res.ComputeSubnet,
	)
	if err != nil {
		// return nil response without any error if no such env found
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return e, nil
}

func (s *state) DeleteEnvironment(ctx context.Context, name string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM environments WHERE name = ?", name)
	if err != nil {
		return err
	}
	return nil
}
