package infrastructure

import (
	"context"
	"fmt"
	"github.com/cloudfauj/cloudfauj/deployment"
	"github.com/hashicorp/terraform-exec/tfexec"
	"os"
	"strings"
	"text/template"
)

const appTfModule = "github.com/cloudfauj/terraform-template.git//app?ref=90fceb4"

const appTfConfig = `module "{{.app}}" {
  source                      = "{{.source}}"
  main_vpc_id                 = module.{{.env}}.main_vpc_id
  ecs_cluster_arn             = module.{{.env}}.compute_ecs_cluster_arn
  compute_subnets             = module.{{.env}}.compute_subnets
  ecs_task_execution_role_arn = module.{{.env}}.ecs_task_execution_role_arn
  env                         = module.{{.env}}.name

  app_name     = "{{.app}}"
  ingress_port = {{.ingress_port}}
  cpu          = {{.cpu}}
  memory       = {{.memory}}
  ecr_image    = "{{.ecr_image}}"
}

output "{{.app}}_ecs_service" {
  value = module.{{.app}}.ecs_service
}`

func (i *Infrastructure) CreateApplication(
	ctx context.Context,
	spec *deployment.Spec,
	tf *tfexec.Terraform,
	tfFile *os.File,
) error {
	appConfig := i.appTfConfig(spec)
	if _, err := tfFile.Write([]byte(appConfig)); err != nil {
		return fmt.Errorf("failed to write Terraform configuration to app file: %v", err)
	}
	if err := tf.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize terraform: %v", err)
	}
	if err := tf.Apply(ctx, tfexec.Target("module."+spec.App.Name)); err != nil {
		return fmt.Errorf("failed to apply terraform changes: %v", err)
	}
	return nil
}

// ModifyApplication is simply a wrapper around CreateApplication because both
// involve writing & applying TF configurations.
func (i *Infrastructure) ModifyApplication(
	ctx context.Context,
	spec *deployment.Spec,
	tf *tfexec.Terraform,
	tfFile *os.File,
) error {
	return i.CreateApplication(ctx, spec, tf, tfFile)
}

func (i *Infrastructure) DestroyApplication(ctx context.Context, tf *tfexec.Terraform, app string) error {
	if err := tf.Destroy(ctx, tfexec.Target("module."+app)); err != nil {
		return fmt.Errorf("failed to destroy terraform infrastructure: %v", err)
	}
	return nil
}

func (i *Infrastructure) AppECSService(ctx context.Context, tf *tfexec.Terraform, app string) (string, error) {
	res, err := tf.Output(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to read terraform output: %v", err)
	}
	service := string(res[app+"_ecs_service"].Value)
	service = strings.Trim(service, "\"")
	return service, nil
}

func (i *Infrastructure) appTfConfig(spec *deployment.Spec) string {
	var b strings.Builder
	t := template.Must(template.New("").Parse(appTfConfig))
	data := map[string]interface{}{
		"env":          spec.TargetEnv,
		"app":          spec.App.Name,
		"source":       appTfModule,
		"ingress_port": spec.App.Resources.Network.BindPort,
		"ecr_image":    spec.Artifact,
		"cpu":          fargateRoundedCPU(spec.App.Resources.Cpu),
		"memory":       fargateRoundedMemory(spec.App.Resources.Cpu, spec.App.Resources.Memory),
	}
	t.Execute(&b, data)
	return b.String()
}
