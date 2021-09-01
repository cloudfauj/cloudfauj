package state

import (
	"context"
	"database/sql"
)

const sqlCreateDomainTable = "CREATE TABLE IF NOT EXISTS domains (name VARCHAR(800) PRIMARY KEY)"

func (s *state) AddDomain(ctx context.Context, name string) error {
	q := "INSERT INTO domains VALUES(?)"
	stmt, err := s.db.PrepareContext(ctx, q)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, name)
	return err
}

func (s *state) CheckDomainExists(ctx context.Context, name string) (bool, error) {
	var res string
	err := s.db.QueryRowContext(ctx, "SELECT name FROM domains WHERE name = ?", name).Scan(&res)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *state) DeleteDomain(ctx context.Context, name string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM domains WHERE name = ?", name)
	return err
}

func (s *state) ListDomains(ctx context.Context) ([]string, error) {
	var res []string

	rows, err := s.db.QueryContext(ctx, "SELECT name FROM domains")
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
