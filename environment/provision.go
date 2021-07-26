package environment

import (
	"context"
	"fmt"
	"github.com/cloudfauj/cloudfauj/infrastructure"
	"strconv"
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

	if err := e.createALBInfra(ctx); err != nil {
		eventsCh <- Event{Err: err}
		return
	}
	eventsCh <- Event{Msg: "Created Load balancer infrastructure"}
}

func (e *Environment) createECSInfra(ctx context.Context) error {
	// create iam role(s)
	role, err := e.Infra.CreateIAMRole(ctx)
	if err != nil {
		return fmt.Errorf("failed to create IAM role for compute: %v", err)
	}
	e.Res.ComputeIAMRole = role

	// create security group
	sg, err := e.Infra.CreateSecurityGroup(ctx)
	if err != nil {
		return fmt.Errorf("failed to create ECS security group: %v", err)
	}
	e.Res.ECSSecurityGroup = sg

	// create ECS fargate cluster
	c, err := e.Infra.CreateFargateCluster(ctx, e.baseResourceName())
	if err != nil {
		return fmt.Errorf("failed to create ECS cluster: %v", err)
	}
	e.Res.ECSCluster = c

	return nil
}

func (e *Environment) createALBInfra(ctx context.Context) error {
	// create subnets
	e.Res.AlbSubnets = make([]string, AlbSubnetCount, AlbSubnetCount)
	az := []string{infrastructure.Az1Suffix, infrastructure.Az2Suffix}

	for i := 0; i < AlbSubnetCount; i++ {
		name := e.baseResourceName() + strconv.Itoa(i+1)
		s, err := e.Infra.CreateSubnet(ctx, name, e.Res.VpcId, az[i], 8)
		if err != nil {
			return fmt.Errorf("failed to create ALB subnet: %v", err)
		}
		e.Res.AlbSubnets[i] = s
	}

	// create security group
	sg, err := e.Infra.CreateSecurityGroup(ctx)
	if err != nil {
		return fmt.Errorf("failed to create ALB security group: %v", err)
	}
	e.Res.AlbSecurityGroup = sg

	// create ALB
	lb, err := e.Infra.CreateALB(ctx, e.baseResourceName(), e.Res.AlbSecurityGroup, e.Res.AlbSubnets)
	if err != nil {
		return fmt.Errorf("failed to create ALB: %v", err)
	}
	e.Res.Alb = lb

	return nil
}

func (e *Environment) baseResourceName() string {
	return "cfoj-" + e.Name
}
