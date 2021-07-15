package deployment

import "github.com/cloudfauj/cloudfauj/application"

// Spec describes the deployment specification supplied by the client
// to deploy a particular application.
type Spec struct {
	App       *application.Application `json:"app"`
	TargetEnv string                   `json:"target_env"`
	Artifact  string                   `json:"artifact"`
}

func (s *Spec) CheckIsValid() error {
	// validate own fields
	return s.App.CheckIsValid()
}
