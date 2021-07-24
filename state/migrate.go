package state

import (
	"context"
	"fmt"
)

const sqlCreateEnvTable = `CREATE TABLE IF NOT EXISTS environments (
	name VARCHAR(100) NOT NULL PRIMARY KEY,
	status VARCHAR(25) NOT NULL,
	vpc_id VARCHAR(40) NOT NULL,
	internet_gateway VARCHAR(50) NOT NULL,
	default_route_table VARCHAR(50) NOT NULL,
	ecs_security_group VARCHAR(100) NOT NULL,
	ecs_cluster VARCHAR(100) NOT NULL,
	fargate_capacity_provider VARCHAR(100) NOT NULL,
	compute_iam_role VARCHAR(200) NOT NULL,
	lb_security_group VARCHAR(100) NOT NULL,
	load_balancer VARCHAR(200) NOT NULL
)`

func (s *state) Migrate(ctx context.Context) error {
	if _, err := s.db.ExecContext(ctx, sqlCreateEnvTable); err != nil {
		return fmt.Errorf("failed to create environments table: %v", err)
	}
	return nil
}
