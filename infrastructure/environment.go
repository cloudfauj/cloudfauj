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

module "{{.env_name}}" {
  source              = "github.com/cloudfauj/terraform-template.git//env?ref=0dcbb03"
  env_name            = "{{.env_name}}"
  main_vpc_cidr_block = "{{.main_vpc_cidr}}"
}

output "ecs_cluster_arn" {
  value = module.{{.env_name}}.compute_ecs_cluster_arn
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
  env_name            = "{{.env_name}}"
  vpc_id              = module.{{.env_name}}.main_vpc_id
  acm_certificate_arn = data.terraform_remote_state.domain.outputs.ssl_cert_arn
  alb_subnets         = module.{{.env_name}}.alb_subnets
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

func (i *Infrastructure) EnvECSCluster(ctx context.Context, tf *tfexec.Terraform) (string, error) {
	res, err := tf.Output(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to read terraform output: %v", err)
	}
	cluster := string(res["ecs_cluster_arn"].Value)
	cluster = strings.Trim(cluster, "\"")
	return cluster, nil
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
		"env_name":            env.Name,
		"domain_tfstate_path": domainTFStateFile,
	}

	t.Execute(&b, data)
	return b.String()
}
