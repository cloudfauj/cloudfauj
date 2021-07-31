package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/sirupsen/logrus"
	"net"
)

const (
	VPCFrozenBits     = 16
	MinVPCCidr        = "10.0.0.0/16"
	LargestVPCCidr    = "10.0.0.0/8"
	ECSTaskExecPolicy = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
)

type Infrastructure struct {
	region string
	log    *logrus.Logger
	ec2    *ec2.Client
	iam    *iam.Client
	ecs    *ecs.Client
	lb     *elasticloadbalancingv2.Client
}

func New(
	l *logrus.Logger,
	ec2 *ec2.Client,
	ecs *ecs.Client,
	i *iam.Client,
	lb *elasticloadbalancingv2.Client,
	region string,
) *Infrastructure {
	return &Infrastructure{
		log: l, ec2: ec2, iam: i, ecs: ecs, lb: lb, region: region,
	}
}

// CreateVPC creates a new VPC in aws with an available CIDR
func (i *Infrastructure) CreateVPC(ctx context.Context) (string, error) {
	c, err := i.nextAvailableCIDR(ctx)
	i.log.Debugf("Using CIDR %s", c)
	if err != nil {
		return "", fmt.Errorf("failed to calculate CIDR: %v", err)
	}
	res, err := i.ec2.CreateVpc(ctx, &ec2.CreateVpcInput{CidrBlock: aws.String(c)})
	if err != nil {
		return "", err
	}
	return aws.ToString(res.Vpc.VpcId), nil
}

// DestroyVPC deletes the given VPC
func (i *Infrastructure) DestroyVPC(ctx context.Context, id string) error {
	_, err := i.ec2.DeleteVpc(ctx, &ec2.DeleteVpcInput{VpcId: aws.String(id)})
	return err
}

// CreateInternetGateway creates a new internet gateway and
// attaches it to the specified VPC.
func (i *Infrastructure) CreateInternetGateway(ctx context.Context, vpc string) (string, error) {
	g, err := i.ec2.CreateInternetGateway(ctx, &ec2.CreateInternetGatewayInput{})
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

// DestroyInternetGateway detaches the internet gateway from its VPC and deletes it.
func (i *Infrastructure) DestroyInternetGateway(ctx context.Context, vpc, gateway string) error {
	gid := aws.String(gateway)
	_, err := i.ec2.DetachInternetGateway(
		ctx,
		&ec2.DetachInternetGatewayInput{VpcId: aws.String(vpc), InternetGatewayId: gid},
	)
	if err != nil {
		return fmt.Errorf("failed to detach internet gateway from vpc: %v", err)
	}
	_, err = i.ec2.DeleteInternetGateway(ctx, &ec2.DeleteInternetGatewayInput{InternetGatewayId: gid})
	return err
}

// AddGatewayRoute adds an internet gateway as the router for all public traffic
// in a VPC's default route table.
func (i *Infrastructure) AddGatewayRoute(ctx context.Context, vpc string, igw string) error {
	rt, err := i.vpcMainRT(ctx, vpc)
	if err != nil {
		return err
	}
	_, err = i.ec2.CreateRoute(ctx, &ec2.CreateRouteInput{
		RouteTableId:         rt.RouteTableId,
		DestinationCidrBlock: aws.String("0.0.0.0/0"),
		GatewayId:            aws.String(igw),
	})
	return err
}

// RemoveGatewayRoute deletes the given route table associated with an internet gateway
func (i *Infrastructure) RemoveGatewayRoute(ctx context.Context, vpc string) error {
	rt, err := i.vpcMainRT(ctx, vpc)
	if err != nil {
		return err
	}
	_, err = i.ec2.DeleteRoute(ctx, &ec2.DeleteRouteInput{
		RouteTableId:         rt.RouteTableId,
		DestinationCidrBlock: aws.String("0.0.0.0/0"),
	})
	return err
}

// CreateSubnet creates a subnet in the given VPC.
func (i *Infrastructure) CreateSubnet(ctx context.Context, name, vpc string, newBits, num int) (string, error) {
	v, err := i.ec2.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{VpcIds: []string{vpc}})
	if err != nil {
		return "", fmt.Errorf("failed to get vpc info: %v", err)
	}
	_, vpcCidr, _ := net.ParseCIDR(aws.ToString(v.Vpcs[0].CidrBlock))

	subCidr, err := cidr.Subnet(vpcCidr, newBits, num)
	if err != nil {
		return "", fmt.Errorf("failed to calculate subnet CIDR: %v", err)
	}
	s, err := i.ec2.CreateSubnet(ctx, &ec2.CreateSubnetInput{
		VpcId:     aws.String(vpc),
		CidrBlock: aws.String(subCidr.String()),
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

// CreateECSTaskExecIAMRole create an IAM role for ECS tasks to assume so they can
// pull images from ECR and ship logs to cloudwatch.
func (i *Infrastructure) CreateECSTaskExecIAMRole(ctx context.Context, name string) (string, error) {
	ad := `{
  "Version": "2012-10-17",
  "Statement": {
    "Effect": "Allow",
    "Principal": {"Service": "ecs-tasks.amazonaws.com"},
    "Action": "sts:AssumeRole"
  }
}`
	r, err := i.iam.CreateRole(ctx, &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(ad),
		RoleName:                 aws.String(name),
		Description:              aws.String("ECS task execution role for Cloudfauj environment"),
	})
	if err != nil {
		return "", err
	}
	_, err = i.iam.AttachRolePolicy(ctx, &iam.AttachRolePolicyInput{
		PolicyArn: aws.String(ECSTaskExecPolicy),
		RoleName:  r.Role.RoleName,
	})
	if err != nil {
		return "", fmt.Errorf("failed to attach policy to role: %v", err)
	}

	p := `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["logs:CreateLogGroup"],
      "Resource": "*"
    }
  ]
}`
	_, err = i.iam.PutRolePolicy(ctx, &iam.PutRolePolicyInput{
		RoleName:       r.Role.RoleName,
		PolicyDocument: aws.String(p),
		PolicyName:     aws.String(name + "-permissions"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to attach inline policy to role: %v", err)
	}
	return aws.ToString(r.Role.Arn), nil
}

// DeleteECSTaskExecIAMRole detaches policies from & deletes an ECS task exec IAM role.
func (i *Infrastructure) DeleteECSTaskExecIAMRole(ctx context.Context, name string) error {
	_, err := i.iam.DetachRolePolicy(ctx, &iam.DetachRolePolicyInput{
		PolicyArn: aws.String(ECSTaskExecPolicy),
		RoleName:  aws.String(name),
	})
	if err != nil {
		return fmt.Errorf("failed to detach policy from role: %v", err)
	}
	_, err = i.iam.DeleteRolePolicy(ctx, &iam.DeleteRolePolicyInput{
		PolicyName: aws.String(name + "-permissions"),
		RoleName:   aws.String(name),
	})
	if err != nil {
		return fmt.Errorf("failed to delete inline policy from role: %v", err)
	}
	_, err = i.iam.DeleteRole(ctx, &iam.DeleteRoleInput{RoleName: aws.String(name)})
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

// nextAvailableCIDR returns the first /16 CIDR available for use in the target AWS account-region
func (i *Infrastructure) nextAvailableCIDR(ctx context.Context) (string, error) {
	// todo: paginate to ensure we have all VPCs
	res, err := i.ec2.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{})
	if err != nil {
		return "", err
	}
	existingCidrs := make([]*net.IPNet, len(res.Vpcs))
	for j, vpc := range res.Vpcs {
		_, ipn, _ := net.ParseCIDR(aws.ToString(vpc.CidrBlock))
		existingCidrs[j] = ipn
	}

	_, super, _ := net.ParseCIDR(LargestVPCCidr)
	_, proposed, _ := net.ParseCIDR(MinVPCCidr)
	for {
		all := append(existingCidrs, proposed)
		if err := cidr.VerifyNoOverlap(all, super); err == nil {
			return proposed.String(), nil
		}
		next, maxed := cidr.NextSubnet(proposed, VPCFrozenBits)
		if maxed || (next.IP[0] > 10) {
			return "", errors.New("no CIDRs available")
		}
		proposed = next
	}
}

// vpcMainRT returns the main route table associated with the given VPC
func (i *Infrastructure) vpcMainRT(ctx context.Context, vpc string) (types.RouteTable, error) {
	rt, err := i.ec2.DescribeRouteTables(ctx, &ec2.DescribeRouteTablesInput{
		Filters: []types.Filter{
			{Name: aws.String("vpc-id"), Values: []string{vpc}},
			{Name: aws.String("association.main"), Values: []string{"true"}},
		},
	})
	if err != nil {
		return types.RouteTable{}, fmt.Errorf("failed to get VPC's main route table: %v", err)
	}
	return rt.RouteTables[0], nil
}
