package infrastructure

import (
	"context"
	"errors"
	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"net"
)

const (
	vpcFrozenBits  = 16
	minVPCCidr     = "10.0.0.0/16"
	largestVPCCidr = "10.0.0.0/8"
)

// nextAvailableCIDR returns the first /16 CIDR available for use in the target AWS account-Region
func (i *Infrastructure) nextAvailableCIDR(ctx context.Context) (string, error) {
	// todo: paginate to ensure we have all VPCs
	res, err := i.Ec2.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{})
	if err != nil {
		return "", err
	}
	existingCidrs := make([]*net.IPNet, len(res.Vpcs))
	for j, vpc := range res.Vpcs {
		_, ipn, _ := net.ParseCIDR(aws.ToString(vpc.CidrBlock))
		existingCidrs[j] = ipn
	}

	_, super, _ := net.ParseCIDR(largestVPCCidr)
	_, proposed, _ := net.ParseCIDR(minVPCCidr)
	for {
		all := append(existingCidrs, proposed)
		if err := cidr.VerifyNoOverlap(all, super); err == nil {
			return proposed.String(), nil
		}
		next, maxed := cidr.NextSubnet(proposed, vpcFrozenBits)
		if maxed || (next.IP[0] > 10) {
			return "", errors.New("no CIDRs available")
		}
		proposed = next
	}
}
