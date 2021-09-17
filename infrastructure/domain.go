package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudfauj/cloudfauj/domain"
	"github.com/hashicorp/terraform-exec/tfexec"
	"strings"
	"text/template"
)

const domainDnsTfConfigTpl = `resource "aws_route53_zone" "dns_manager" {
  name          = "{{.domain_name}}"
  tags          = local.common_tags
  force_destroy = true
  comment       = "Public Hosted Zone for {{.domain_name}} managed by Cloudfauj"
}

// Add the DNS validation records to the Hosted zone.
// Note that this alone is not enough to validate the ACM cert.
// The NS records of the main R53 zone must be configured with the domain
// provider manually by the user.
// After that, ACM will be able to validate & issue the certificate.
resource "aws_route53_record" "acm_cert_validation" {
  for_each = {
    for dvo in aws_acm_certificate.primary_wildcard_cert.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  allow_overwrite = true

  name    = each.value.name
  records = [each.value.record]
  ttl     = 60
  type    = each.value.type
  zone_id = aws_route53_zone.dns_manager.zone_id
}

output "zone_id" {
  value = aws_route53_zone.dns_manager.zone_id
}

output "name_servers" {
  value = aws_route53_zone.dns_manager.name_servers
}`

const domainCertTfConfigTpl = `resource "aws_acm_certificate" "primary_wildcard_cert" {
  domain_name               = "{{.domain_name}}"
  subject_alternative_names = ["*.{{.domain_name}}"]
  validation_method         = "DNS"
  tags                      = local.common_tags
}

output "ssl_cert_arn" {
  value = aws_acm_certificate.primary_wildcard_cert.arn
}

output "apex_domain" {
  value = "{{.domain_name}}"
}`

// DomainTFConfig returns a map.
// The keys are names of terraform config files needed for the domain infrastructure.
// The values are their corresponding TF code.
// This method generates the TF configuration depending on the components being used
// for the domain.
func (i *Infrastructure) DomainTFConfig(d *domain.Domain) (map[string]string, error) {
	// NOTE: As of now, only route53 dns service & acm cert authority are supported,
	// so this method generates tf only for those, regardless of what's specified
	// in the domain configuration.
	res := map[string]string{
		"terraform.tf":      i.tfCoreConfig(),
		"dns_service.tf":    i.domainTfConfig(d, domainDnsTfConfigTpl),
		"cert_authority.tf": i.domainTfConfig(d, domainCertTfConfigTpl),
	}
	return res, nil
}

// CreateDomain creates infrastructure for a domain.
// It returns the Name Server records of the DNS hosted zone.
func (i *Infrastructure) CreateDomain(ctx context.Context, tf *tfexec.Terraform) ([]string, error) {
	var nsRecords []string

	if err := tf.Init(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize terraform: %v", err)
	}
	if err := tf.Apply(ctx); err != nil {
		return nil, fmt.Errorf("failed to apply TF config: %v", err)
	}

	// read NS records from terraform output
	res, err := tf.Output(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read terraform output: %v", err)
	}
	if err := json.Unmarshal(res["name_servers"].Value, &nsRecords); err != nil {
		return nil, fmt.Errorf("failed to parse terraform output: %v", err)
	}
	return nsRecords, nil
}

// DeleteDomain destroys the Terraform infrastructure of a domain
func (i *Infrastructure) DeleteDomain(ctx context.Context, tf *tfexec.Terraform) error {
	return tf.Destroy(ctx)
}

func (i *Infrastructure) domainTfConfig(d *domain.Domain, tpl string) string {
	var b strings.Builder

	t := template.Must(template.New("").Parse(tpl))
	data := map[string]interface{}{"domain_name": d.Name}

	t.Execute(&b, data)
	return b.String()
}
