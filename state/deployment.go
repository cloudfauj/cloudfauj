package state

import (
	"context"
	"github.com/cloudfauj/cloudfauj/deployment"
)

func (s *state) GetDeployment(ctx context.Context, id string) (*deployment.Deployment, error) {
	// return nil if depoyment doesn't exist
	return nil, nil
}

func (s *state) ListDeployments(ctx context.Context, statuses ...string) ([]*deployment.Deployment, error) {
	return nil, nil
}
