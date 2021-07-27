package deployment

import (
	"errors"
	"github.com/cloudfauj/cloudfauj/application"
	"strings"
)

// Spec describes the deployment specification supplied by the client
// to deploy a particular application.
type Spec struct {
	App       *application.Application `json:"app"`
	TargetEnv string                   `json:"target_env"`
	Artifact  string                   `json:"artifact"`
}

func (s *Spec) CheckIsValid() error {
	if len(strings.TrimSpace(s.TargetEnv)) == 0 {
		return errors.New("target environment not specified")
	}
	if len(strings.TrimSpace(s.Artifact)) == 0 {
		return errors.New("artifact not specified")
	}
	return s.App.CheckIsValid()
}
