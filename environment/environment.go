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
	Name   string     `json:"name"`
	Status string     `json:"status"`
	Res    *Resources `json:"resources"`
	Infra  *infrastructure.Infrastructure
}

type Resources struct {
	VpcId             string `json:"vpc_id"`
	InternetGateway   string `json:"internet_gateway"`
	DefaultRouteTable string `json:"default_route_table"`

	ECSCluster      string `json:"ecs_cluster"`
	TaskExecIAMRole string `json:"task_exec_iam_role"`
	ComputeSubnet   string `json:"compute_subnet"`
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
