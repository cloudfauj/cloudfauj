package environment

import (
	"errors"
	"strings"
)

const (
	StatusProvisioning = "provisioning"
	StatusProvisioned  = "provisioned"
	StatusDestroying   = "destroying"
)

// Environment represents a Cloudfauj environment
type Environment struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

func (e *Environment) CheckIsValid() error {
	// TODO: Ensure env name doesn't use any of the reserved names
	//  (eg: any that start with _)
	if len(strings.TrimSpace(e.Name)) == 0 {
		return errors.New("name cannot be empty")
	}
	return nil
}
