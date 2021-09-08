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

const envTfConfigTpl = `{{.tf_core_config}}

module "main_env" {
  source              = "github.com/cloudfauj/terraform-template.git//env?ref=0dcbb03"
  env_name            = "{{.env_name}}"
  main_vpc_cidr_block = "{{.main_vpc_cidr}}"
}

output "ecs_cluster_arn" {
  value = module.main_env.compute_ecs_cluster_arn
}

output "main_vpc_id" {
  value = module.main_env.main_vpc_id
}

output "compute_subnets" {
  value = module.main_env.compute_subnets
}

output "ecs_task_execution_role_arn" {
  value = module.main_env.ecs_task_execution_role_arn
}

output "name" {
  value = module.main_env.name
}

{{.alb_module}}`

const envAlbTfConfigTpl = `data "terraform_remote_state" "domain" {
  backend = "local"

  config = {
    path = "{{.domain_tfstate_path}}"
  }
}

module "apps_load_balancer" {
  source              = "github.com/cloudfauj/terraform-template.git//env/load_balancer?ref=0dcbb03"
  env_name            = module.main_env.name
  vpc_id              = module.main_env.main_vpc_id
  acm_certificate_arn = data.terraform_remote_state.domain.outputs.ssl_cert_arn
  alb_subnets         = module.main_env.alb_subnets
}

output "apps_alb_name" {
  value = module.apps_load_balancer.apps_alb_name
}

output "main_alb_https_listener" {
  value = module.apps_load_balancer.main_alb_https_listener
}`

func (i *Infrastructure) CreateEnvironment(
	ctx context.Context,
	env *environment.Environment,
	domainTfStateFile string,
	tf *tfexec.Terraform,
	tfFile *os.File,
) error {
	cidr, err := i.nextAvailableCIDR(ctx)
	if err != nil {
		return fmt.Errorf("failed to compute VPC CIDR: %v", err)
	}

	envConf := i.envTfConfig(env, cidr, domainTfStateFile)
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

func (i *Infrastructure) envTfConfig(env *environment.Environment, vpcCidr, domainTFStateFile string) string {
	var b strings.Builder
	t := template.Must(template.New("").Parse(envTfConfigTpl))
	data := map[string]interface{}{
		"tf_core_config": i.tfCoreConfig(),
		"env_name":       env.Name,
		"main_vpc_cidr":  vpcCidr,
		"alb_module":     i.envAlbTfConfig(env, domainTFStateFile),
	}

	t.Execute(&b, data)
	return b.String()
}

func (i *Infrastructure) envAlbTfConfig(env *environment.Environment, domainTFStateFile string) string {
	if !env.DomainEnabled() {
		return ""
	}

	var b strings.Builder
	t := template.Must(template.New("").Parse(envAlbTfConfigTpl))
	data := map[string]interface{}{
		"domain_tfstate_path": domainTFStateFile,
	}

	t.Execute(&b, data)
	return b.String()
}
