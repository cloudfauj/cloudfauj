package infrastructure

import "context"

type Infrastructure struct{}

func New() *Infrastructure {
	return &Infrastructure{}
}

func (i *Infrastructure) GetAvailableCIDR(ctx context.Context, frozenBits int) (string, error) {
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
