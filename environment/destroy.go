package environment

import (
	"context"
	"fmt"
)

func (e *Environment) Destroy(ctx context.Context, eventsCh chan<- Event) {
	defer close(eventsCh)

	if err := e.destroyALBInfra(ctx); err != nil {
		eventsCh <- Event{Err: err}
		return
	}
	eventsCh <- Event{Msg: "Destroyed Load balancer infrastructure"}

	if err := e.destroyECSInfra(ctx); err != nil {
		eventsCh <- Event{Err: err}
		return
	}
	eventsCh <- Event{Msg: "Destroyed ECS fargate infrastructure"}

	// destroy public route table
	if err := e.Infra.DestroyPublicRouteTable(ctx, e.Res.DefaultRouteTable); err != nil {
		eventsCh <- Event{Err: fmt.Errorf("failed to destroy default route table: %v", err)}
		return
	}
	eventsCh <- Event{Msg: "Destroyed default route table"}

	// destroy inet gateway
	if err := e.Infra.DestroyInternetGateway(ctx, e.Res.InternetGateway); err != nil {
		eventsCh <- Event{Err: fmt.Errorf("failed to destroy internet gateway: %v", err)}
		return
	}
	eventsCh <- Event{Msg: "Destroyed internet gateway"}

	// destroy VPC
	if err := e.Infra.DestroyVPC(ctx, e.Res.VpcId); err != nil {
		eventsCh <- Event{Err: fmt.Errorf("failed to destroy VPC: %v", err)}
		return
	}
	eventsCh <- Event{Msg: "Destroyed VPC"}
}

func (e *Environment) destroyALBInfra(ctx context.Context) error {
	// destroy security group
	if err := e.Infra.DestroySecurityGroup(ctx, e.Res.AlbSecurityGroup); err != nil {
		return fmt.Errorf("failed to destroy ALB security group: %v", err)
	}

	// destroy ALB
	if err := e.Infra.DestroyALB(ctx, e.Res.Alb); err != nil {
		return fmt.Errorf("failed to destroy ALB: %v", err)
	}

	return nil
}

func (e *Environment) destroyECSInfra(ctx context.Context) error {
	// destroy iam role(s)
	if err := e.Infra.DeleteIAMRole(ctx, e.Res.ComputeIAMRole); err != nil {
		return fmt.Errorf("failed to destroy IAM role for compute: %v", err)
	}

	// destroy security group
	if err := e.Infra.DestroySecurityGroup(ctx, e.Res.ECSSecurityGroup); err != nil {
		return fmt.Errorf("failed to destroy ECS security group: %v", err)
	}

	// destroy ECS fargate cluster
	if err := e.Infra.DestroyFargateCluster(ctx, e.Res.ECSCluster); err != nil {
		return fmt.Errorf("failed to destroy ECS cluster: %v", err)
	}

	return nil
}
