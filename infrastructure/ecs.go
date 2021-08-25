package infrastructure

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"strconv"
)

func (i *Infrastructure) ECSService(ctx context.Context, service, cluster string) (types.Service, error) {
	res, err := i.Ecs.DescribeServices(ctx, &ecs.DescribeServicesInput{
		Services: []string{service},
		Cluster:  aws.String(cluster),
	})
	if err != nil {
		return types.Service{}, err
	}
	return res.Services[0], nil
}

func (i *Infrastructure) ECSServiceStatus(ctx context.Context, service, cluster string) (string, error) {
	s, err := i.ECSService(ctx, service, cluster)
	if err != nil {
		return "", err
	}
	return aws.ToString(s.Status), nil
}

func (i *Infrastructure) ECSServicePrimaryDeployment(ctx context.Context, service, cluster string) (types.Deployment, error) {
	s, err := i.ECSService(ctx, service, cluster)
	if err != nil {
		return types.Deployment{}, err
	}
	// todo: ensure that the first item in Deployments list is always the PRIMARY deployment
	return s.Deployments[0], nil
}

// memRange returns discrete memory values (MB) from start to end
// at increments of 1024.
func memRange(start, end int) []int {
	var res []int
	inc := 1024
	for i := start; i <= end; i += inc {
		res = append(res, i)
	}
	return res
}

// fargateRoundedCPU returns the amount of CPU compatible with fargate.
// It is at least as much as the user-specified CPU.
func fargateRoundedCPU(cpu int) string {
	rng := []int{0, 256, 512, 1024, 2048, 4096}
	for i := 0; i < len(rng)-1; i++ {
		if cpu > rng[i] && cpu <= rng[i+1] {
			return strconv.Itoa(rng[i+1])
		}
	}
	// todo: return err if cpu > max rng in fargate
	return strconv.Itoa(rng[len(rng)-1])
}

// fargateRoundedMemory returns the amount of Memory compatible with fargate.
// It is at least as much as the user-specified memory.
func fargateRoundedMemory(cpu, memory int) string {
	ranges := map[string][]int{
		"256":  {512, 1024, 2048},
		"512":  memRange(1024, 4096),
		"1024": memRange(2048, 8192),
		"2048": memRange(4096, 16384),
		"4096": memRange(9216, 30720),
	}
	rng := ranges[fargateRoundedCPU(cpu)]
	if memory <= rng[0] {
		return strconv.Itoa(rng[0])
	}
	for i := 0; i < len(rng)-1; i++ {
		if memory <= rng[i+1] {
			return strconv.Itoa(rng[i+1])
		}
	}
	// todo: return err if memory > max rng in fargate
	return strconv.Itoa(rng[len(rng)-1])
}
