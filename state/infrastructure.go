package state

import (
	"context"
	"github.com/cloudfauj/cloudfauj/infrastructure"
)

const sqlCreateAppInfraTable = `CREATE TABLE IF NOT EXISTS app_infra (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	app VARCHAR(100) NOT NULL,
	env VARCHAR(100) NOT NULL,
	ecs_task_definition VARCHAR(300) NOT NULL,
	ecs_service VARCHAR(300) NOT NULL,
	security_group VARCHAR(300) NOT NULL,
	UNIQUE(app, env)
)`

func (s *state) CreateAppInfra(ctx context.Context, infra *infrastructure.AppInfra) error {
	q := `INSERT INTO app_infra(
	app, env, ecs_task_definition, ecs_service, security_group
) VALUES(?, ?, ?, ?, ?)`

	stmt, err := s.db.PrepareContext(ctx, q)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(
		ctx,
		infra.App,
		infra.Env,
		infra.EcsTaskDefinition,
		infra.ECSService,
		infra.SecurityGroup,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *state) UpdateAppInfra(ctx context.Context, infra *infrastructure.AppInfra) error {
	q := `UPDATE app_infra
SET
	ecs_task_definition = ?,
	ecs_service = ?,
	security_group = ?
WHERE app = ? AND env = ?`

	stmt, err := s.db.PrepareContext(ctx, q)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(
		ctx,
		infra.EcsTaskDefinition,
		infra.ECSService,
		infra.SecurityGroup,
		infra.App,
		infra.Env,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *state) AppInfra(ctx context.Context, name, env string) (*infrastructure.AppInfra, error) {
	var (
		id int
		i  infrastructure.AppInfra
	)
	err := s.db.QueryRowContext(
		ctx, "SELECT * FROM app_infra WHERE app = ? AND env = ?", name, env,
	).Scan(
		&id,
		&i.App,
		&i.Env,
		&i.EcsTaskDefinition,
		&i.ECSService,
		&i.SecurityGroup,
	)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func (s *state) DeleteAppInfra(ctx context.Context, name, env string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM app_infra WHERE app = ? AND env = ?", name, env)
	return err
}
