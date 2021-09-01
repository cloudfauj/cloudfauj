package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-exec/tfexec"
	"os"
	"strings"
	"text/template"
)

const domainModuleSource = "github.com/cloudfauj/terraform-template.git//domain?ref=1b2921c"

const domainTfConfigTpl = `{{.tf_core_config}}

module "{{.module_name}}" {
  source = "{{.module_source}}"
  name   = "{{.domain_name}}"
}

output "name_servers" {
  value = module.{{.module_name}}.name_servers
}

output "zone_id" {
  value = module.{{.module_name}}.zone_id
}

output "ssl_cert_arn" {
  value = module.{{.module_name}}.ssl_cert_arn
}`

// CreateDomain creates infrastructure to manage a domain.
// It returns the Name Server records of the DNS hosted zone.
func (i *Infrastructure) CreateDomain(
	ctx context.Context, name string, tf *tfexec.Terraform, tfFile *os.File,
) ([]string, error) {
	var nsRecords []string

	conf := i.domainTFConfig(name)
	if _, err := tfFile.Write([]byte(conf)); err != nil {
		return nil, fmt.Errorf("failed to write Terraform configuration to file: %v", err)
	}
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

func (i *Infrastructure) DeleteDomain(ctx context.Context, tf *tfexec.Terraform) error {
	return tf.Destroy(ctx)
}

func (i *Infrastructure) domainTFConfig(name string) string {
	var b strings.Builder

	t := template.Must(template.New("").Parse(domainTfConfigTpl))
	data := map[string]interface{}{
		"tf_core_config": i.tfCoreConfig(),
		"module_name":    domainModuleName(name),
		"module_source":  domainModuleSource,
		"domain_name":    name,
	}

	t.Execute(&b, data)
	return b.String()
}

func domainModuleName(name string) string {
	return strings.ReplaceAll(name, ".", "_")
}
