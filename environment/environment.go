package environment

import "context"

const (
	StatusProvisioning = "provisioning"
	StatusProvisioned  = "provisioned"
)

type Environment struct {
	Name   string `json:"name"`
	Status string
}

type Event struct {
	Msg string
	Err error
}

func (e *Environment) CheckIsValid() error {
	// mandatory fields to not be empty
	// format should be correct (regex)
	return nil
}

func (e *Environment) Provision(ctx context.Context, eventsCh chan<- Event) {
	close(eventsCh)
}
