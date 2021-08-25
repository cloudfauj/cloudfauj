package infrastructure

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/sirupsen/logrus"
)

// Interacts with AWS to provision and manage infrastructure resources.
type Infrastructure struct {
	Log      *logrus.Logger
	Region   string
	Ec2      *ec2.Client
	Ecs      *ecs.Client
	TFBinary string
}
