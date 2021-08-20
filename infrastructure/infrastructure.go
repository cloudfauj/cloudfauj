package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/sirupsen/logrus"
	"net"
	"os"
	"path"
	"strings"
	"text/template"
)

const (
	VPCFrozenBits               = 16
	MinVPCCidr                  = "10.0.0.0/16"
	LargestVPCCidr              = "10.0.0.0/8"
	TerraformAwsProviderVersion = "3.54.0"
)

const tfConfigTpl = `terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "{{.aws_provider_version}}"
    }
  }
}

provider "aws" {
  region = "{{.aws_region}}"
}`

type Infrastructure struct {
	region       string
	log          *logrus.Logger
	ec2          *ec2.Client
	ecs          *ecs.Client
	tfConfigDir  string
	tfBinaryPath string
}

func New(
	l *logrus.Logger,
	ec2 *ec2.Client,
	ecs *ecs.Client,
	region, tfDir, tfBin string,
) *Infrastructure {
	return &Infrastructure{
		log:          l,
		ec2:          ec2,
		ecs:          ecs,
		region:       region,
		tfConfigDir:  tfDir,
		tfBinaryPath: tfBin,
	}
}

// NextAvailableCIDR returns the first /16 CIDR available for use in the target AWS account-region
func (i *Infrastructure) NextAvailableCIDR(ctx context.Context) (string, error) {
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

func (i *Infrastructure) Tf(workSubDir string) (*tfexec.Terraform, error) {
	tf, err := tfexec.NewTerraform(path.Join(i.tfConfigDir, workSubDir), i.tfBinaryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create new terraform object: %s", err)
	}
	// Pass the server process' environment variables to Terraform process
	tf.SetEnv(nil)
	// Set logging
	tf.SetLogger(i.log)
	tf.SetStderr(os.Stderr)
	tf.SetStdout(os.Stdout)

	return tf, nil
}

func (i *Infrastructure) TfConfig() string {
	var b strings.Builder
	t := template.Must(template.New("").Parse(tfConfigTpl))
	data := map[string]interface{}{
		"aws_region":           i.region,
		"aws_provider_version": TerraformAwsProviderVersion,
	}
	t.Execute(&b, data)
	return b.String()
}
