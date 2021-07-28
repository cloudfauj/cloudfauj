package environment

import (
	"context"
	"fmt"
)

func (e *Environment) Provision(ctx context.Context, eventsCh chan<- Event) {
	defer close(eventsCh)

	v, err := e.Infra.CreateVPC(ctx)
	if err != nil {
		eventsCh <- Event{Err: fmt.Errorf("failed to create VPC: %v", err)}
		return
	}
	e.Res.VpcId = v
	eventsCh <- Event{Msg: "Created VPC"}

	g, err := e.Infra.CreateInternetGateway(ctx, e.Res.VpcId)
	if err != nil {
		eventsCh <- Event{Err: fmt.Errorf("failed to create internet gateway: %v", err)}
		return
	}
	e.Res.InternetGateway = g
	eventsCh <- Event{Msg: "Created Internet Gateway"}

	if err := e.Infra.AddGatewayRoute(ctx, e.Res.VpcId, e.Res.InternetGateway); err != nil {
		eventsCh <- Event{Err: fmt.Errorf("failed to create default route table: %v", err)}
		return
	}
	eventsCh <- Event{Msg: "Setup routing"}

	if err := e.createECSInfra(ctx); err != nil {
		eventsCh <- Event{Err: err}
		return
	}
	eventsCh <- Event{Msg: "Created ECS Fargate infrastructure"}
}

func (e *Environment) createECSInfra(ctx context.Context) error {
	s, err := e.Infra.CreateSubnet(ctx, e.baseResourceName(), e.Res.VpcId, 4, 1)
	if err != nil {
		return fmt.Errorf("failed to create subnet for containers: %v", err)
	}
	e.Res.ComputeSubnet = s

	// create ECS task execution IAM role that allows tasks
	// to pull images & ship logs to CWL.
	n := e.baseResourceName() + "-ecs-task-exec"
	if _, err := e.Infra.CreateECSTaskExecIAMRole(ctx, n); err != nil {
		return fmt.Errorf("failed to create IAM role for compute: %v", err)
	}
	e.Res.TaskExecIAMRole = n

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
