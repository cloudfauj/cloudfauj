// Code generated by smithy-go-codegen DO NOT EDIT.

package endpoints

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/endpoints"
	"regexp"
)

// Options is the endpoint resolver configuration options
type Options struct {
	DisableHTTPS bool
}

// Resolver IAM endpoint resolver
type Resolver struct {
	partitions endpoints.Partitions
}

// ResolveEndpoint resolves the service endpoint for the given region and options
func (r *Resolver) ResolveEndpoint(region string, options Options) (endpoint aws.Endpoint, err error) {
	if len(region) == 0 {
		return endpoint, &aws.MissingRegionError{}
	}

	opt := endpoints.Options{
		DisableHTTPS: options.DisableHTTPS,
	}
	return r.partitions.ResolveEndpoint(region, opt)
}

// New returns a new Resolver
func New() *Resolver {
	return &Resolver{
		partitions: defaultPartitions,
	}
}

var partitionRegexp = struct {
	Aws      *regexp.Regexp
	AwsCn    *regexp.Regexp
	AwsIso   *regexp.Regexp
	AwsIsoB  *regexp.Regexp
	AwsUsGov *regexp.Regexp
}{

	Aws:      regexp.MustCompile("^(us|eu|ap|sa|ca|me|af)\\-\\w+\\-\\d+$"),
	AwsCn:    regexp.MustCompile("^cn\\-\\w+\\-\\d+$"),
	AwsIso:   regexp.MustCompile("^us\\-iso\\-\\w+\\-\\d+$"),
	AwsIsoB:  regexp.MustCompile("^us\\-isob\\-\\w+\\-\\d+$"),
	AwsUsGov: regexp.MustCompile("^us\\-gov\\-\\w+\\-\\d+$"),
}

var defaultPartitions = endpoints.Partitions{
	{
		ID: "aws",
		Defaults: endpoints.Endpoint{
			Hostname:          "iam.{region}.amazonaws.com",
			Protocols:         []string{"https"},
			SignatureVersions: []string{"v4"},
		},
		RegionRegex:       partitionRegexp.Aws,
		IsRegionalized:    false,
		PartitionEndpoint: "aws-global",
		Endpoints: endpoints.Endpoints{
			"aws-global": endpoints.Endpoint{
				Hostname: "iam.amazonaws.com",
				CredentialScope: endpoints.CredentialScope{
					Region: "us-east-1",
				},
			},
			"iam-fips": endpoints.Endpoint{
				Hostname: "iam-fips.amazonaws.com",
				CredentialScope: endpoints.CredentialScope{
					Region: "us-east-1",
				},
			},
		},
	},
	{
		ID: "aws-cn",
		Defaults: endpoints.Endpoint{
			Hostname:          "iam.{region}.amazonaws.com.cn",
			Protocols:         []string{"https"},
			SignatureVersions: []string{"v4"},
		},
		RegionRegex:       partitionRegexp.AwsCn,
		IsRegionalized:    false,
		PartitionEndpoint: "aws-cn-global",
		Endpoints: endpoints.Endpoints{
			"aws-cn-global": endpoints.Endpoint{
				Hostname: "iam.cn-north-1.amazonaws.com.cn",
				CredentialScope: endpoints.CredentialScope{
					Region: "cn-north-1",
				},
			},
		},
	},
	{
		ID: "aws-iso",
		Defaults: endpoints.Endpoint{
			Hostname:          "iam.{region}.c2s.ic.gov",
			Protocols:         []string{"https"},
			SignatureVersions: []string{"v4"},
		},
		RegionRegex:       partitionRegexp.AwsIso,
		IsRegionalized:    false,
		PartitionEndpoint: "aws-iso-global",
		Endpoints: endpoints.Endpoints{
			"aws-iso-global": endpoints.Endpoint{
				Hostname: "iam.us-iso-east-1.c2s.ic.gov",
				CredentialScope: endpoints.CredentialScope{
					Region: "us-iso-east-1",
				},
			},
		},
	},
	{
		ID: "aws-iso-b",
		Defaults: endpoints.Endpoint{
			Hostname:          "iam.{region}.sc2s.sgov.gov",
			Protocols:         []string{"https"},
			SignatureVersions: []string{"v4"},
		},
		RegionRegex:       partitionRegexp.AwsIsoB,
		IsRegionalized:    false,
		PartitionEndpoint: "aws-iso-b-global",
		Endpoints: endpoints.Endpoints{
			"aws-iso-b-global": endpoints.Endpoint{
				Hostname: "iam.us-isob-east-1.sc2s.sgov.gov",
				CredentialScope: endpoints.CredentialScope{
					Region: "us-isob-east-1",
				},
			},
		},
	},
	{
		ID: "aws-us-gov",
		Defaults: endpoints.Endpoint{
			Hostname:          "iam.{region}.amazonaws.com",
			Protocols:         []string{"https"},
			SignatureVersions: []string{"v4"},
		},
		RegionRegex:       partitionRegexp.AwsUsGov,
		IsRegionalized:    false,
		PartitionEndpoint: "aws-us-gov-global",
		Endpoints: endpoints.Endpoints{
			"aws-us-gov-global": endpoints.Endpoint{
				Hostname: "iam.us-gov.amazonaws.com",
				CredentialScope: endpoints.CredentialScope{
					Region: "us-gov-west-1",
				},
			},
			"iam-govcloud-fips": endpoints.Endpoint{
				Hostname: "iam.us-gov.amazonaws.com",
				CredentialScope: endpoints.CredentialScope{
					Region: "us-gov-west-1",
				},
			},
		},
	},
}
