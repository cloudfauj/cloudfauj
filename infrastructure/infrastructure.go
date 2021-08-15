package infrastructure

import (
	"context"
	"errors"
	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/sirupsen/logrus"
	"net"
)

const (
	VPCFrozenBits  = 16
	MinVPCCidr     = "10.0.0.0/16"
	LargestVPCCidr = "10.0.0.0/8"
)

type Infrastructure struct {
	Tf     *tfexec.Terraform
	region string
	log    *logrus.Logger
	ec2    *ec2.Client
	ecs    *ecs.Client
}

func New(
	l *logrus.Logger,
	tf *tfexec.Terraform,
	ec2 *ec2.Client,
	ecs *ecs.Client,
	region string,
) *Infrastructure {
	return &Infrastructure{log: l, Tf: tf, ec2: ec2, ecs: ecs, region: region}
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
