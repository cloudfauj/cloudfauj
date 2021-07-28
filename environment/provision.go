package environment

import (
	"context"
	"fmt"
	"github.com/cloudfauj/cloudfauj/infrastructure"
)

func (e *Environment) Provision(ctx context.Context, eventsCh chan<- Event) {
	defer close(eventsCh)

	// create VPC
	v, err := e.Infra.CreateVPC(ctx)
	if err != nil {
		eventsCh <- Event{Err: fmt.Errorf("failed to create VPC: %v", err)}
		return
	}
	e.Res.VpcId = v
	eventsCh <- Event{Msg: "Created VPC"}

	// create internet gateway for VPC
	g, err := e.Infra.CreateInternetGateway(ctx, e.Res.VpcId)
	if err != nil {
		eventsCh <- Event{Err: fmt.Errorf("failed to create internet gateway: %v", err)}
		return
	}
	e.Res.InternetGateway = g
	eventsCh <- Event{Msg: "Created Internet Gateway"}

	// create default route table
	rt, err := e.Infra.CreatePublicRouteTable(ctx, e.Res.VpcId, e.Res.InternetGateway)
	if err != nil {
		eventsCh <- Event{Err: fmt.Errorf("failed to create default route table: %v", err)}
		return
	}
	e.Res.DefaultRouteTable = rt
	eventsCh <- Event{Msg: "Created default route table"}

	if err := e.createECSInfra(ctx); err != nil {
		eventsCh <- Event{Err: err}
		return
	}
	eventsCh <- Event{Msg: "Created ECS Fargate infrastructure"}
}

func (e *Environment) createECSInfra(ctx context.Context) error {
	n := e.baseResourceName() + infrastructure.Az1Suffix
	s, err := e.Infra.CreateSubnet(ctx, n, e.Res.VpcId, infrastructure.Az1Suffix, 4)
	if err != nil {
		return fmt.Errorf("failed to create subnet for containers: %v", err)
	}
	e.Res.ComputeSubnet = s

	// create ECS task execution IAM role that allows tasks
	// to pull images & ship logs to CWL.
	role, err := e.Infra.CreateIAMRole(ctx)
	if err != nil {
		return fmt.Errorf("failed to create IAM role for compute: %v", err)
	}
	e.Res.TaskExecIAMRole = role

	// create ECS fargate cluster
	c, err := e.Infra.CreateFargateCluster(ctx, e.baseResourceName())
	if err != nil {
		return fmt.Errorf("failed to create ECS cluster: %v", err)
	}
	e.Res.ECSCluster = c

	return nil
}

func (e *Environment) baseResourceName() string {
	return "cfoj-" + e.Name
}
