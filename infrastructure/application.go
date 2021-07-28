package infrastructure

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type AppInfra struct {
	App               string `json:"app"`
	EcsTaskDefinition string `json:"ecs_task_definition"`
	TargetGroup       string `json:"target_group"`
	AlbListenerRule   string `json:"alb_listener_rule"`
	DNSRecord         string `json:"dns_record"`
	ECSService        string `json:"ecs_service"`
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

func (i *Infrastructure) CreateSecurityGroup(ctx context.Context) (string, error) {
	return "", nil
}

func (i *Infrastructure) DestroySecurityGroup(ctx context.Context, id string) error {
	_, err := i.ec2.DeleteSecurityGroup(ctx, &ec2.DeleteSecurityGroupInput{GroupId: aws.String(id)})
	return err
}
