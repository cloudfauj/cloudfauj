package state

import (
	"context"
	"database/sql"
	"github.com/cloudfauj/cloudfauj/application"
)

const sqlCreateAppTable = `CREATE TABLE IF NOT EXISTS applications (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name VARCHAR(100) NOT NULL,
	env VARCHAR(100) NOT NULL,
	type VARCHAR(40) NOT NULL,
	visibility VARCHAR(40) NOT NULL,
	health_path VARCHAR(70) NOT NULL,
	cpu INT NOT NULL,
	memory INT NOT NULL,
	bind_port INT NOT NULL,
	UNIQUE(name, env)
)`

func (s *state) CreateApp(ctx context.Context, app *application.Application, env string) error {
	q := `INSERT INTO applications(
	name, env, type, visibility, health_path, cpu, memory, bind_port
) VALUES(?, ?, ?, ?, ?, ?, ?, ?)`

	stmt, err := s.db.PrepareContext(ctx, q)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(
		ctx,
		app.Name,
		env,
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

func (s *state) UpdateApp(ctx context.Context, app *application.Application, env string) error {
	q := `UPDATE applications
SET
	type = ?,
	visibility = ?,
	health_path = ?,
	cpu = ?,
	memory = ?,
	bind_port = ?
WHERE name = ? AND env = ?`

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
		env,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *state) App(ctx context.Context, name, env string) (*application.Application, error) {
	var (
		id int
		e  string
	)
	a := &application.Application{
		HealthCheck: &application.HealthCheck{},
		Resources:   &application.Resources{Network: &application.Network{}},
	}
	err := s.db.QueryRowContext(
		ctx, "SELECT * FROM applications WHERE name = ? AND env = ?", name, env,
	).Scan(
		&id,
		&a.Name,
		&e,
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

func (s *state) DeleteApp(ctx context.Context, name, env string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM applications WHERE name = ? AND env = ?", name, env)
	return err
}

//func (s *state) CheckAppExists(ctx context.Context, name, env string) (bool, error) {
//	var res string
//	err := s.db.QueryRowContext(
//		ctx, "SELECT name FROM applications WHERE name = ? AND env = ?", name, env,
//	).Scan(&res)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			return false, nil
//		}
//		return false, err
//	}
//	return true, nil
//}
