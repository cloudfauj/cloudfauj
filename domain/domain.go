package domain

import (
	"errors"
	"strings"
)

const DNSServiceRoute53 = "aws_route53"
const CertAuthorityACM = "aws_acm"

type Domain struct {
	Name                 string `json:"name"`
	DNSService           string `json:"dns_service" mapstructure:"dns_service"`
	CertificateAuthority string `json:"cert_authority" mapstructure:"cert_authority"`
}

func (d *Domain) CheckIsValid() error {
	if len(strings.TrimSpace(d.Name)) == 0 {
		return errors.New("name cannot be empty")
	}
	if d.DNSService != DNSServiceRoute53 {
		return errors.New("only " + DNSServiceRoute53 + " DNS service is supported as of now")
	}
	if d.CertificateAuthority != CertAuthorityACM {
		return errors.New("only " + CertAuthorityACM + " certificate authority is supported as of now")
	}
	return nil
}
