package state

import (
	"context"
	"github.com/cloudfauj/cloudfauj/infrastructure"
)

func (s *state) CreateAppInfra(context.Context, *infrastructure.AppInfra) error {
	return nil
}

func (s *state) UpdateAppInfra(context.Context, *infrastructure.AppInfra) error {
	return nil
}

func (s *state) AppInfra(context.Context, string) (*infrastructure.AppInfra, error) {
	return nil, nil
}
