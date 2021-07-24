package state

import (
	"context"
	"fmt"
)

func (s *state) Migrate(ctx context.Context) error {
	if _, err := s.db.ExecContext(ctx, sqlCreateEnvTable); err != nil {
		return fmt.Errorf("failed to create environments table: %v", err)
	}
	return nil
}
