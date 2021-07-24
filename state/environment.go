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
	default_route_table VARCHAR(50),
	ecs_security_group VARCHAR(100),
	ecs_cluster VARCHAR(100),
	fargate_capacity_provider VARCHAR(100),
	compute_iam_role VARCHAR(200),
	lb_security_group VARCHAR(100),
	load_balancer VARCHAR(200)
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
	default_route_table,
	ecs_security_group,
	ecs_cluster,
	fargate_capacity_provider,
	compute_iam_role,
	lb_security_group,
	load_balancer
) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

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
		e.Res.DefaultRouteTable,
		e.Res.ECSSecurityGroup,
		e.Res.ECSCluster,
		e.Res.FargateCapProvider,
		e.Res.ComputeIAMRole,
		e.Res.AlbSecurityGroup,
		e.Res.Alb,
	)
	if err != nil {
		return err
	}
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
