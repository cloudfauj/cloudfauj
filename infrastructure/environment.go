package infrastructure

import (
	"context"
	"fmt"
	"github.com/cloudfauj/cloudfauj/environment"
	"github.com/hashicorp/terraform-exec/tfexec"
	"os"
	"strings"
	"text/template"
)

const envModuleSource = "github.com/cloudfauj/terraform-template.git//env?ref=90fceb4"

const envTfConfigTpl = `{{.tf_core_config}}

module "{{.name}}" {
  source              = "{{.module_source}}"
  env_name            = "{{.name}}"
  main_vpc_cidr_block = "{{.main_vpc_cidr}}"
}

output "ecs_cluster_arn" {
  value = module.{{.name}}.compute_ecs_cluster_arn
}`

func (i *Infrastructure) CreateEnvironment(
	ctx context.Context,
	env *environment.Environment,
	tf *tfexec.Terraform,
	tfFile *os.File,
) error {
	cidr, err := i.NextAvailableCIDR(ctx)
	if err != nil {
		return fmt.Errorf("failed to compute VPC CIDR: %v", err)
	}

	envConf := i.envTfConfig(env.Name, cidr)
	if _, err := tfFile.Write([]byte(envConf)); err != nil {
		return fmt.Errorf("failed to write Terraform configuration to file: %v", err)
	}
	if err := tf.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize terraform: %v", err)
	}
	if err := tf.Apply(ctx); err != nil {
		return fmt.Errorf("failed to apply TF config: %v", err)
	}
	return nil
}

func (i *Infrastructure) DestroyEnvironment(ctx context.Context, tf *tfexec.Terraform) error {
	// since the working directory only contains TF configuration of the env,
	// the whole config can be safely deleted without impacting any infra outside
	// the env.
	if err := tf.Destroy(ctx); err != nil {
		return fmt.Errorf("failed to destroy: %v", err)
	}
	return nil
}

func (i *Infrastructure) envTfConfig(env, vpcCidr string) string {
	var b strings.Builder
	t := template.Must(template.New("").Parse(envTfConfigTpl))
	data := map[string]interface{}{
		"tf_core_config": i.tfCoreConfig(),
		"name":           env,
		"module_source":  envModuleSource,
		"main_vpc_cidr":  vpcCidr,
	}
	t.Execute(&b, data)
	return b.String()
}
