package environment

import (
	"errors"
	"github.com/cloudfauj/cloudfauj/infrastructure"
	"strings"
)

const AlbSubnetCount = 2

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

	ECSSecurityGroup string `json:"ecs_security_group"`
	ECSCluster       string `json:"ecs_cluster"`
	ComputeIAMRole   string `json:"compute_iam_role"`

	AlbSecurityGroup string   `json:"lb_security_group"`
	AlbSubnets       []string `json:"lb_subnets"`
	Alb              string   `json:"load_balancer"`
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
