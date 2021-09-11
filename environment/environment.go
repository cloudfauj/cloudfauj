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

// A Cloudfauj Environment that can contain applications
type Environment struct {
	Name   string `json:"name"`
	Domain string `json:"domain"`
	Status string `json:"status"`
}

func (e *Environment) CheckIsValid() error {
	// TODO
	//  1. Ensure env name doesn't use any of the reserved names (eg: any that start with _)
	//  2. Any validations/blacklists to be applied to domain?
	if len(strings.TrimSpace(e.Name)) == 0 {
		return errors.New("name cannot be empty")
	}
	return nil
}

// DomainEnabled returns true if a domain is associated with the environment
func (e *Environment) DomainEnabled() bool {
	return len(strings.TrimSpace(e.Domain)) != 0
}
