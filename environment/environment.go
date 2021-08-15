package environment

import (
	"errors"
	"github.com/cloudfauj/cloudfauj/infrastructure"
	"strings"
)

const (
	StatusProvisioning = "provisioning"
	StatusProvisioned  = "provisioned"
	StatusDestroying   = "destroying"
)

type Environment struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Infra  *infrastructure.Infrastructure
}

type Event struct {
	Msg string
	Err error
}

func (e *Environment) CheckIsValid() error {
	if len(strings.TrimSpace(e.Name)) == 0 {
		return errors.New("name cannot be empty")
	}
	return nil
}
