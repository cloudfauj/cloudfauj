package state

import (
	"context"
	"fmt"
)

func (s *state) Migrate(ctx context.Context) error {
	if _, err := s.db.ExecContext(ctx, sqlCreateEnvTable); err != nil {
		return fmt.Errorf("failed to create environments table: %v", err)
	}
	if _, err := s.db.ExecContext(ctx, sqlCreateAppTable); err != nil {
		return fmt.Errorf("failed to create applications table: %v", err)
	}
	if _, err := s.db.ExecContext(ctx, sqlCreateAppInfraTable); err != nil {
		return fmt.Errorf("failed to create applications infra table: %v", err)
	}
	return nil
}
