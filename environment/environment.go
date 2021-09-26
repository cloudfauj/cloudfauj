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

const NetworkAWS = "aws"
const OrchFargate = "aws_ecs_fargate"
const LoadBalALB = "aws_alb"

// A Cloudfauj Environment that can contain applications
type Environment struct {
	Name         string `json:"name"`
	Network      string `json:"network"`
	Orchestrator string `json:"orchestrator"`

	Domain       string `json:"domain"`
	LoadBalancer string `json:"load_balancer" mapstructure:"load_balancer"`

	Status string `json:"status"`
}

func (e *Environment) CheckIsValid() error {
	// TODO
	//  1. Ensure env name doesn't use any of the reserved names (eg: any that start with _)
	//  2. Any validations/blacklists to be applied to domain?
	if len(strings.TrimSpace(e.Name)) == 0 {
		return errors.New("name cannot be empty")
	}
	if e.Network != NetworkAWS {
		return errors.New("only " + NetworkAWS + " network type is supported for now")
	}
	if e.Orchestrator != OrchFargate {
		return errors.New("only " + OrchFargate + " container orchestrator is supported for now")
	}
	if e.DomainEnabled() && e.LoadBalancer != LoadBalALB {
		return errors.New("only " + LoadBalALB + " load balancer is supported for now")
	}
	return nil
}

// DomainEnabled returns true if a domain is associated with the environment
func (e *Environment) DomainEnabled() bool {
	return len(strings.TrimSpace(e.Domain)) != 0
}
