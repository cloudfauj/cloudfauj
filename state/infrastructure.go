package state

import (
	"context"
	"github.com/cloudfauj/cloudfauj/infrastructure"
)

const sqlCreateAppInfraTable = `CREATE TABLE IF NOT EXISTS app_infra (
	app VARCHAR(100) NOT NULL PRIMARY KEY,
	ecs_task_definition VARCHAR(300) NOT NULL,
	ecs_service VARCHAR(300) NOT NULL,
	security_group VARCHAR(300) NOT NULL
)`

func (s *state) CreateAppInfra(ctx context.Context, infra *infrastructure.AppInfra) error {
	q := `INSERT INTO app_infra(
	app, ecs_task_definition, ecs_service, security_group
) VALUES(?, ?, ?, ?)`

	stmt, err := s.db.PrepareContext(ctx, q)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(
		ctx,
		infra.App,
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
WHERE app = ?`

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
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *state) AppInfra(ctx context.Context, name string) (*infrastructure.AppInfra, error) {
	var i infrastructure.AppInfra
	err := s.db.QueryRowContext(
		ctx, "SELECT * FROM app_infra WHERE app = ?", name,
	).Scan(
		&i.App,
		&i.EcsTaskDefinition,
		&i.ECSService,
		&i.SecurityGroup,
	)
	if err != nil {
		return nil, err
	}
	return &i, nil
}
