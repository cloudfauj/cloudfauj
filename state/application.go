package state

import (
	"context"
	"database/sql"
	"github.com/cloudfauj/cloudfauj/application"
)

const sqlCreateAppTable = `CREATE TABLE IF NOT EXISTS applications (
	name VARCHAR(100) NOT NULL PRIMARY KEY,
	type VARCHAR(40) NOT NULL,
	visibility VARCHAR(40) NOT NULL,
	health_path VARCHAR(70) NOT NULL,
	cpu INT NOT NULL,
	memory INT NOT NULL,
	bind_port INT NOT NULL
)`

func (s *state) CreateApp(ctx context.Context, app *application.Application) error {
	q := `INSERT INTO applications(
	name, type, visibility, health_path, cpu, memory, bind_port
) VALUES(?, ?, ?, ?, ?, ?, ?)`

	stmt, err := s.db.PrepareContext(ctx, q)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(
		ctx,
		app.Name,
		app.Type,
		app.Visibility,
		app.HealthCheck.Path,
		app.Resources.Cpu,
		app.Resources.Memory,
		app.Resources.Network.BindPort,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *state) CheckAppExists(ctx context.Context, name string) (bool, error) {
	var res string
	err := s.db.QueryRowContext(ctx, "SELECT name FROM applications WHERE name = ?", name).Scan(&res)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *state) UpdateApp(ctx context.Context, app *application.Application) error {
	q := `UPDATE applications
SET
	type = ?,
	visibility = ?,
	health_path = ?,
	cpu = ?,
	memory = ?,
	bind_port = ?
WHERE name = ?`

	stmt, err := s.db.PrepareContext(ctx, q)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(
		ctx,
		app.Type,
		app.Visibility,
		app.HealthCheck.Path,
		app.Resources.Cpu,
		app.Resources.Memory,
		app.Resources.Network.BindPort,
		app.Name,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *state) App(ctx context.Context, name string) (*application.Application, error) {
	a := &application.Application{
		HealthCheck: &application.HealthCheck{},
		Resources:   &application.Resources{Network: &application.Network{}},
	}
	err := s.db.QueryRowContext(
		ctx, "SELECT * FROM applications WHERE name = ?", name,
	).Scan(
		&a.Name,
		&a.Type,
		&a.Visibility,
		&a.HealthCheck.Path,
		&a.Resources.Cpu,
		&a.Resources.Memory,
		&a.Resources.Network.BindPort,
	)
	if err != nil {
		// return nil response without any error if no such app found
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return a, nil
}

func (s *state) DeleteApp(ctx context.Context, name string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM applications WHERE name = ?", name)
	return err
}
