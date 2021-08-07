package infrastructure

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	types2 "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"strconv"
)

type AppInfra struct {
	App               string `json:"app"`
	Env               string `json:"env"`
	EcsTaskDefinition string `json:"ecs_task_definition"`
	ECSService        string `json:"ecs_service"`
	SecurityGroup     string `json:"security_group"`
}

// TaskDefintionParams contains the parameters supplied to CreateTaskDefinition()
type TaskDefintionParams struct {
	Env          string
	Service      string
	TaskExecRole string
	Image        string
	Cpu          int
	Memory       int
	BindPort     int32
}

// RoundedCPU returns the amount of CPU compatible with fargate.
// It is at least as much as the user-specified CPU.
func (p *TaskDefintionParams) RoundedCPU() string {
	rng := []int{0, 256, 512, 1024, 2048, 4096}
	for i := 0; i < len(rng)-1; i++ {
		if p.Cpu > rng[i] && p.Cpu <= rng[i+1] {
			return strconv.Itoa(rng[i+1])
		}
	}
	// todo: return err if cpu > max rng in fargate
	return strconv.Itoa(rng[len(rng)-1])
}

// RoundedMemory returns the amount of Memory compatible with fargate.
// It is at least as much as the user-specified memory.
func (p *TaskDefintionParams) RoundedMemory() string {
	ranges := map[string][]int{
		"256":  {512, 1024, 2048},
		"512":  memRange(1024, 4096),
		"1024": memRange(2048, 8192),
		"2048": memRange(4096, 16384),
		"4096": memRange(9216, 30720),
	}
	rng := ranges[p.RoundedCPU()]
	if p.Memory <= rng[0] {
		return strconv.Itoa(rng[0])
	}
	for i := 0; i < len(rng)-1; i++ {
		if p.Memory <= rng[i+1] {
			return strconv.Itoa(rng[i+1])
		}
	}
	// todo: return err if memory > max rng in fargate
	return strconv.Itoa(rng[len(rng)-1])
}

// ECSServiceParams contains the parameters supplied to CreateTaskDefinition()
type ECSServiceParams struct {
	Env           string
	Service       string
	Cluster       string
	TaskDef       string
	ComputeSubnet string
	SecurityGroup string
}

// CreateTaskDefinition creates a new Fargate-compatible task definition.
// If a task def with same family name already exists, this method creates
// a new revision of it.
func (i *Infrastructure) CreateTaskDefinition(ctx context.Context, p *TaskDefintionParams) (string, error) {
	ctr := types.ContainerDefinition{
		Name:  aws.String(p.Service),
		Image: aws.String(p.Image),
		LogConfiguration: &types.LogConfiguration{
			LogDriver: types.LogDriverAwslogs,
			Options: map[string]string{
				"awslogs-create-group":  "true",
				"awslogs-region":        i.region,
				"awslogs-group":         p.Env,
				"awslogs-stream-prefix": p.Service,
			},
		},
		PortMappings: []types.PortMapping{
			{ContainerPort: aws.Int32(p.BindPort)},
		},
		Essential: aws.Bool(true),
	}
	td, err := i.ecs.RegisterTaskDefinition(ctx, &ecs.RegisterTaskDefinitionInput{
		Family:                  aws.String(p.Env + "-" + p.Service),
		ContainerDefinitions:    []types.ContainerDefinition{ctr},
		RequiresCompatibilities: []types.Compatibility{types.CompatibilityFargate},
		ExecutionRoleArn:        aws.String(p.TaskExecRole),
		Cpu:                     aws.String(p.RoundedCPU()),
		Memory:                  aws.String(p.RoundedMemory()),
		NetworkMode:             types.NetworkModeAwsvpc,
	})
	if err != nil {
		return "", err
	}
	return aws.ToString(td.TaskDefinition.TaskDefinitionArn), nil
}

// CreateSecurityGroup creates a security group with TCP ingress allowed from a given port.
func (i *Infrastructure) CreateSecurityGroup(ctx context.Context, env, service, vpc string, port int32) (string, error) {
	sg, err := i.ec2.CreateSecurityGroup(ctx, &ec2.CreateSecurityGroupInput{
		GroupName:   aws.String(env + "-" + service),
		VpcId:       aws.String(vpc),
		Description: aws.String(fmt.Sprintf("Traffic control for %s/%s", env, service)),
		TagSpecifications: []types2.TagSpecification{{
			ResourceType: types2.ResourceTypeSecurityGroup,
			Tags:         []types2.Tag{{Key: aws.String("Name"), Value: aws.String(env + "-" + service)}},
		}},
	})
	if err != nil {
		return "", err
	}
	_, err = i.ec2.AuthorizeSecurityGroupIngress(ctx, &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: sg.GroupId,
		IpPermissions: []types2.IpPermission{
			{
				FromPort: aws.Int32(port),
				ToPort:   aws.Int32(port),
				IpRanges: []types2.IpRange{
					{
						CidrIp:      aws.String("0.0.0.0/0"),
						Description: aws.String("Application ingress traffic"),
					},
				},
				IpProtocol: aws.String("tcp"),
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to authorize ingress: %v", err)
	}
	return aws.ToString(sg.GroupId), nil
}

func (i *Infrastructure) DestroySecurityGroup(ctx context.Context, id string) error {
	_, err := i.ec2.DeleteSecurityGroup(ctx, &ec2.DeleteSecurityGroupInput{GroupId: aws.String(id)})
	return err
}

func (i *Infrastructure) CreateECSService(ctx context.Context, p *ECSServiceParams) (string, error) {
	s, err := i.ecs.CreateService(ctx, &ecs.CreateServiceInput{
		ServiceName:  aws.String(p.Service),
		Cluster:      aws.String(p.Cluster),
		DesiredCount: aws.Int32(1),
		LaunchType:   types.LaunchTypeFargate,
		DeploymentConfiguration: &types.DeploymentConfiguration{
			DeploymentCircuitBreaker: &types.DeploymentCircuitBreaker{Enable: true, Rollback: true},
			MaximumPercent:           aws.Int32(200),
			MinimumHealthyPercent:    aws.Int32(100),
		},
		NetworkConfiguration: &types.NetworkConfiguration{
			AwsvpcConfiguration: &types.AwsVpcConfiguration{
				Subnets:        []string{p.ComputeSubnet},
				AssignPublicIp: types.AssignPublicIpEnabled,
				SecurityGroups: []string{p.SecurityGroup},
			},
		},
		SchedulingStrategy: types.SchedulingStrategyReplica,
		TaskDefinition:     aws.String(p.TaskDef),
	})
	if err != nil {
		return "", err
	}
	return aws.ToString(s.Service.ServiceArn), nil
}

func (i *Infrastructure) ECSService(ctx context.Context, service, cluster string) (types.Service, error) {
	res, err := i.ecs.DescribeServices(ctx, &ecs.DescribeServicesInput{
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

// DestroyECSService destroys an ECS service.
// It drains the service out if any tasks are still running as part of it.
func (i *Infrastructure) DestroyECSService(ctx context.Context, service, cluster string) error {
	_, err := i.ecs.DeleteService(ctx, &ecs.DeleteServiceInput{
		Service: aws.String(service),
		Cluster: aws.String(cluster),
		Force:   aws.Bool(true),
	})
	return err
}

// UpdateECSService updates an ECS service with a task definition, therefore triggering
// a deployment in it.
func (i *Infrastructure) UpdateECSService(ctx context.Context, s, c, t string) error {
	_, err := i.ecs.UpdateService(ctx, &ecs.UpdateServiceInput{
		Service:        aws.String(s),
		Cluster:        aws.String(c),
		TaskDefinition: aws.String(t),
	})
	return err
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
