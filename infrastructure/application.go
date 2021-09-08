package infrastructure

import (
	"context"
	"fmt"
	"github.com/cloudfauj/cloudfauj/deployment"
	"github.com/cloudfauj/cloudfauj/environment"
	"github.com/hashicorp/terraform-exec/tfexec"
	"os"
	"strings"
	"text/template"
)

const appTfConfig = `{{.tf_core_config}}

data "terraform_remote_state" "main_env" {
  backend = "local"

  config = {
    path = "{{.env_tfstate_path}}"
  }
}

# Variables that need to be supplied during invokation
# Note that these have default empty values only to make TF destroy
# invokation easier.
variable "ingress_port" { default = 0 }
variable "cpu" { default = 256 }
variable "memory" { default = 512 }
variable "ecr_image" { default = "" }
variable "app_health_check_path" { default = "" }

module "main_app" {
  source                      = "github.com/cloudfauj/terraform-template.git//app?ref=f32a060"
  main_vpc_id                 = data.terraform_remote_state.main_env.outputs.main_vpc_id
  ecs_cluster_arn             = data.terraform_remote_state.main_env.outputs.ecs_cluster_arn
  compute_subnets             = data.terraform_remote_state.main_env.outputs.compute_subnets
  ecs_task_execution_role_arn = data.terraform_remote_state.main_env.outputs.ecs_task_execution_role_arn
  env                         = data.terraform_remote_state.main_env.outputs.name

  app_name             = "{{.app}}"
  lb_target_group_arns = [{{.app_target_group_arn}}]

  # Modifiable values are read in from variables
  ingress_port = var.ingress_port
  cpu          = var.cpu
  memory       = var.memory
  ecr_image    = var.ecr_image
}

output "ecs_cluster_arn" {
  value = data.terraform_remote_state.main_env.outputs.ecs_cluster_arn
}

output "ecs_service" {
  value = module.main_app.ecs_service
}

{{.app_domain_hook_module}}`

const appDomainHookTfConfigTpl = `data "terraform_remote_state" "domain" {
  backend = "local"

  config = {
    path = "{{.domain_tfstate_path}}"
  }
}

module "app_domain_hook" {
  source                 = "github.com/cloudfauj/terraform-template.git//app/domain_hook?ref=81178d6"
  app_name               = "{{.app}}"
  app_health_check_path  = var.app_health_check_path
  env_name               = data.terraform_remote_state.main_env.outputs.name
  env_vpc_id             = data.terraform_remote_state.main_env.outputs.main_vpc_id
  env_apps_alb_name      = data.terraform_remote_state.main_env.outputs.apps_alb_name
  alb_listener_https_arn = data.terraform_remote_state.main_env.outputs.main_alb_https_listener
  apex_domain            = data.terraform_remote_state.domain.outputs.apex_domain
  route53_zone_id        = data.terraform_remote_state.domain.outputs.zone_id
}`

// Objects supplied to the CreateApplication method
type CreateApplicationInput struct {
	// Deployment specification of the application
	Spec *deployment.Spec

	// Target environment
	Env *environment.Environment

	// Terraform object with working dir set to the application's
	Tf *tfexec.Terraform

	// The file to write the application's terraform configuration to
	TfFile *os.File

	// If target environment has domain enabled, the exact path on the system
	// of the file containing the domain's TF state.
	DomainTFStateFile string

	// The exact path on the system of the file containing the environment's
	// TF state.
	EnvTFStateFile string
}

func (i *Infrastructure) CreateApplication(ctx context.Context, input *CreateApplicationInput) error {
	appConfig := i.appTfConfig(input)
	if _, err := input.TfFile.Write([]byte(appConfig)); err != nil {
		return fmt.Errorf("failed to write Terraform configuration to app file: %v", err)
	}
	if err := input.Tf.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize terraform: %v", err)
	}
	if err := i.applyAppConfig(ctx, input.Spec, input.Tf); err != nil {
		return fmt.Errorf("failed to apply terraform changes: %v", err)
	}
	return nil
}

func (i *Infrastructure) ModifyApplication(
	ctx context.Context,
	spec *deployment.Spec,
	tf *tfexec.Terraform,
) error {
	return i.applyAppConfig(ctx, spec, tf)
}

func (i *Infrastructure) applyAppConfig(ctx context.Context, spec *deployment.Spec, tf *tfexec.Terraform) error {
	return tf.Apply(
		ctx,
		tfexec.Var("app_health_check_path="+spec.App.HealthCheck.Path),
		tfexec.Var(fmt.Sprintf("cpu=%v", fargateRoundedCPU(spec.App.Resources.Cpu))),
		tfexec.Var(
			fmt.Sprintf(
				"memory=%v", fargateRoundedMemory(spec.App.Resources.Cpu, spec.App.Resources.Memory),
			),
		),
		tfexec.Var(fmt.Sprintf("ingress_port=%d", spec.App.Resources.Network.BindPort)),
		tfexec.Var("ecr_image="+spec.Artifact),
	)
}

func (i *Infrastructure) DestroyApplication(ctx context.Context, tf *tfexec.Terraform) error {
	if err := tf.Destroy(ctx); err != nil {
		return fmt.Errorf("failed to destroy app infrastructure: %v", err)
	}
	return nil
}

func (i *Infrastructure) AppECSService(ctx context.Context, tf *tfexec.Terraform) (string, error) {
	return i.tfOutput(ctx, tf, "ecs_service")
}

func (i *Infrastructure) AppECSCluster(ctx context.Context, tf *tfexec.Terraform) (string, error) {
	return i.tfOutput(ctx, tf, "ecs_cluster_arn")
}

func (i *Infrastructure) tfOutput(ctx context.Context, tf *tfexec.Terraform, varName string) (string, error) {
	res, err := tf.Output(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to read terraform output: %v", err)
	}
	value := string(res[varName].Value)
	value = strings.Trim(value, "\"")
	return value, nil
}

func (i *Infrastructure) appTfConfig(input *CreateApplicationInput) string {
	var b strings.Builder

	tgARN := ""
	if input.Env.DomainEnabled() {
		tgARN = "module.app_domain_hook.target_group_arn"
	}

	t := template.Must(template.New("").Parse(appTfConfig))
	data := map[string]interface{}{
		"app":                    input.Spec.App.Name,
		"tf_core_config":         i.tfCoreConfig(),
		"app_target_group_arn":   tgARN,
		"app_domain_hook_module": i.appDomainHookTfConfig(input),
		"env_tfstate_path":       input.EnvTFStateFile,
	}

	t.Execute(&b, data)
	return b.String()
}

func (i *Infrastructure) appDomainHookTfConfig(input *CreateApplicationInput) string {
	if !input.Env.DomainEnabled() {
		return ""
	}

	var b strings.Builder
	t := template.Must(template.New("").Parse(appDomainHookTfConfigTpl))
	data := map[string]interface{}{
		"app":                 input.Spec.App.Name,
		"domain_tfstate_path": input.DomainTFStateFile,
	}

	t.Execute(&b, data)
	return b.String()
}
