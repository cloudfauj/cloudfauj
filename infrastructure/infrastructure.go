package infrastructure

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/sirupsen/logrus"
)

const (
	Az1Suffix = "a"
	Az2Suffix = "b"
	Az3Suffix = "c"
)

type Infrastructure struct {
	log *logrus.Logger
	ec2 *ec2.Client
	iam *iam.Client
	ecs *ecs.Client
	lb  *elasticloadbalancingv2.Client
}

func New(
	l *logrus.Logger,
	ec2 *ec2.Client,
	ecs *ecs.Client,
	i *iam.Client,
	lb *elasticloadbalancingv2.Client,
) *Infrastructure {
	return &Infrastructure{log: l, ec2: ec2, iam: i, ecs: ecs, lb: lb}
}

// CreateVPC creates a new VPC in aws with an available /16 CIDR
func (i *Infrastructure) CreateVPC(ctx context.Context) (string, error) {
	// todo: find available cidr
	//res, err := i.ec2.CreateVpc(ctx, &ec2.CreateVpcInput{CidrBlock: nil})
	//if err != nil {
	//	return "", err
	//}
	//return aws.ToString(res.Vpc.VpcId), nil

	return "", nil
}

// DestroyVPC deletes the given VPC
func (i *Infrastructure) DestroyVPC(ctx context.Context, id string) error {
	_, err := i.ec2.DeleteVpc(ctx, &ec2.DeleteVpcInput{VpcId: aws.String(id)})
	return err
}

// CreateSubnet creates a subnet in the given VPC-AZ.
// It calculates & uses the next available CIDR based on specified frozen bits.
// eg- if frozen bits = 4 & VPC is /16, then it uses the next /20 subnet
// (16 + 4 frozen) available in the VPC.
func (i *Infrastructure) CreateSubnet(ctx context.Context, name, vpc, azSuffix string, frozen int) (string, error) {
	// todo: calculate CIDR
	// todo: infer AZ

	s, err := i.ec2.CreateSubnet(ctx, &ec2.CreateSubnetInput{
		CidrBlock:        aws.String(cidr),
		VpcId:            aws.String(vpc),
		AvailabilityZone: aws.String(az),
	})
	if err != nil {
		return "", err
	}
	return aws.ToString(s.Subnet.SubnetId), nil
}

// DestroySubnet deletes the given subnet
func (i *Infrastructure) DestroySubnet(ctx context.Context, id string) error {
	_, err := i.ec2.DeleteSubnet(ctx, &ec2.DeleteSubnetInput{SubnetId: aws.String(id)})
	return err
}

// CreateFargateCluster creates an ECS cluster with a default
// provider strategy of Fargate.
func (i *Infrastructure) CreateFargateCluster(ctx context.Context, name string) (string, error) {
	c, err := i.ecs.CreateCluster(ctx, &ecs.CreateClusterInput{
		CapacityProviders: []string{"FARGATE"},
		ClusterName:       aws.String(name),
	})
	if err != nil {
		return "", err
	}
	return aws.ToString(c.Cluster.ClusterArn), nil
}

// DestroyFargateCluster deletes a cluster which only contains Fargate capacity provider
func (i *Infrastructure) DestroyFargateCluster(ctx context.Context, arn string) error {
	_, err := i.ecs.DeleteCluster(ctx, &ecs.DeleteClusterInput{Cluster: aws.String(arn)})
	return err
}

// CreateALB creates an application load balancer
func (i *Infrastructure) CreateALB(ctx context.Context, name, sg string, subnets []string) (string, error) {
	alb, err := i.lb.CreateLoadBalancer(ctx, &elasticloadbalancingv2.CreateLoadBalancerInput{
		Name:           aws.String(name),
		Scheme:         "internet-facing",
		SecurityGroups: []string{sg},
		Subnets:        subnets,
		Type:           "application",
	})
	if err != nil {
		return "", err
	}
	return aws.ToString(alb.LoadBalancers[0].LoadBalancerArn), nil
}

func (i *Infrastructure) DestroyALB(ctx context.Context, arn string) error {
	_, err := i.lb.DeleteLoadBalancer(
		ctx,
		&elasticloadbalancingv2.DeleteLoadBalancerInput{LoadBalancerArn: aws.String(arn)},
	)
	return err
}

func (i *Infrastructure) CreateSecurityGroup(ctx context.Context) (string, error) {
	return "", nil
}

func (i *Infrastructure) DestroySecurityGroup(ctx context.Context, id string) error {
	_, err := i.ec2.DeleteSecurityGroup(ctx, &ec2.DeleteSecurityGroupInput{GroupId: aws.String(id)})
	return err
}

func (i *Infrastructure) CreateIAMRole(ctx context.Context) (string, error) {
	return "", nil
}

func (i *Infrastructure) DeleteIAMRole(ctx context.Context, name string) error {
	_, err := i.iam.DeleteRole(ctx, &iam.DeleteRoleInput{RoleName: aws.String(name)})
	return err
}

// CreateInternetGateway creates a new internet gateway and
// attaches it to the specified VPC.
func (i *Infrastructure) CreateInternetGateway(ctx context.Context, vpc string) (string, error) {
	g, err := i.ec2.CreateInternetGateway(ctx, &ec2.CreateInternetGatewayInput{TagSpecifications: nil})
	if err != nil {
		return "", err
	}

	gid := g.InternetGateway.InternetGatewayId
	_, err = i.ec2.AttachInternetGateway(
		ctx,
		&ec2.AttachInternetGatewayInput{InternetGatewayId: gid, VpcId: aws.String(vpc)},
	)
	if err != nil {
		return "", fmt.Errorf("failed to attach to vpc: %v", err)
	}

	return aws.ToString(gid), nil
}

func (i *Infrastructure) DestroyInternetGateway(ctx context.Context, id string) error {
	_, err := i.ec2.DeleteInternetGateway(ctx, &ec2.DeleteInternetGatewayInput{InternetGatewayId: aws.String(id)})
	return err
}

// CreatePublicRouteTable creates a public route table for a vpc
// by routing all traffic via its internet gateway.
func (i *Infrastructure) CreatePublicRouteTable(ctx context.Context, vpc string, igw string) (string, error) {
	rt, err := i.ec2.CreateRouteTable(ctx, &ec2.CreateRouteTableInput{VpcId: aws.String(vpc)})
	if err != nil {
		return "", err
	}
	_, err = i.ec2.CreateRoute(ctx, &ec2.CreateRouteInput{
		RouteTableId:         rt.RouteTable.RouteTableId,
		DestinationCidrBlock: aws.String("0.0.0.0/0"),
		GatewayId:            aws.String(igw),
	})
	if err != nil {
		return "", err
	}
	return aws.ToString(rt.RouteTable.RouteTableId), nil
}

// DestroyPublicRouteTable deletes the given route table associated with an internet gateway
func (i *Infrastructure) DestroyPublicRouteTable(ctx context.Context, id string) error {
	_, err := i.ec2.DeleteRouteTable(ctx, &ec2.DeleteRouteTableInput{RouteTableId: aws.String(id)})
	return err
}
