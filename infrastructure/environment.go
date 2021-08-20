package infrastructure

import (
	"strings"
	"text/template"
)

const EnvModuleSource = "github.com/cloudfauj/terraform-template.git//env?ref=90fceb4"

const envTfConfigTpl = `{{.tf_core_config}}

module "{{.name}}" {
  source              = "{{.module_source}}"
  env_name            = "{{.name}}"
  main_vpc_cidr_block = "{{.main_vpc_cidr}}"
}

output "ecs_cluster_arn" {
  value = module.{{.name}}.compute_ecs_cluster_arn
}`

func (i *Infrastructure) EnvTfConfig(env, vpcCidr string) string {
	var b strings.Builder
	t := template.Must(template.New("").Parse(envTfConfigTpl))
	data := map[string]interface{}{
		"tf_core_config": i.TfConfig(),
		"name":           env,
		"module_source":  EnvModuleSource,
		"main_vpc_cidr":  vpcCidr,
	}
	t.Execute(&b, data)
	return b.String()
}
