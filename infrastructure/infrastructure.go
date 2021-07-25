package infrastructure

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/sirupsen/logrus"
)

type Infrastructure struct {
	log *logrus.Logger
	ec2 *ec2.Client
	iam *iam.Client
	ecs *ecs.Client
}

func New(l *logrus.Logger, ec2 *ec2.Client, ecs *ecs.Client, i *iam.Client) *Infrastructure {
	return &Infrastructure{log: l, ec2: ec2, iam: i, ecs: ecs}
}

func (i *Infrastructure) GetAvailableCIDR(ctx context.Context, frozenBits int) (string, error) {
	//res, err := i.ec2.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{VpcIds: []string{""}})
	//i.log.Info(err)
	//i.log.Info(res)
	return "", nil
}

func (i *Infrastructure) CreateVPC(ctx context.Context, cidr string) (string, error) {
	return "", nil
}

func (i *Infrastructure) CreateECSCluster(ctx context.Context) (string, error) {
	return "", nil
}

func (i *Infrastructure) CreateALB(ctx context.Context) (string, error) {
	return "", nil
}

func (i *Infrastructure) CreateSecurityGroup(ctx context.Context) (string, error) {
	return "", nil
}

func (i *Infrastructure) CreateFargateCapacityProvider(ctx context.Context) (string, error) {
	return "", nil
}

func (i *Infrastructure) CreateIAMRole(ctx context.Context) (string, error) {
	return "", nil
}

func (i *Infrastructure) CreateInternetGateway(ctx context.Context, vpc string) (string, error) {
	return "", nil
}

func (i *Infrastructure) CreatePublicRouteTable(ctx context.Context, vpc string) (string, error) {
	return "", nil
}

func (i *Infrastructure) CreateTaskDefinition(ctx context.Context) (string, error) {
	return "", nil
}

func (i *Infrastructure) CreateTargetGroup(ctx context.Context) (string, error) {
	return "", nil
}

func (i *Infrastructure) AttachTargetGroup(ctx context.Context, t string) (string, error) {
	return "", nil
}

func (i *Infrastructure) CreateDNSRecord(ctx context.Context) (string, error) {
	return "", nil
}

func (i *Infrastructure) CreateECSService(ctx context.Context) (string, error) {
	return "", nil
}

func (i *Infrastructure) UpdateECSService(ctx context.Context, t string) error {
	return nil
}

func (i *Infrastructure) DestroyVPC(ctx context.Context, id string) error {
	return nil
}

func (i *Infrastructure) DestroyECSCluster(ctx context.Context, id string) error {
	return nil
}

func (i *Infrastructure) DestroyALB(ctx context.Context, id string) error {
	return nil
}

func (i *Infrastructure) DestroySecurityGroup(ctx context.Context, id string) error {
	return nil
}

func (i *Infrastructure) DestroyFargateCapacityProvider(ctx context.Context, id string) error {
	return nil
}

func (i *Infrastructure) DestroyIAMRole(ctx context.Context, id string) error {
	return nil
}

func (i *Infrastructure) DestroyInternetGateway(ctx context.Context, id string) error {
	return nil
}

func (i *Infrastructure) DestroyPublicRouteTable(ctx context.Context, id string) error {
	return nil
}
