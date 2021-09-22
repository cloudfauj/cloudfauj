package infrastructure

import (
	"context"
	"fmt"
	"github.com/cloudfauj/cloudfauj/deployment"
	"github.com/cloudfauj/cloudfauj/environment"
	"github.com/hashicorp/terraform-exec/tfexec"
	"strings"
	"text/template"
)

// A set of Objects supplied to the AppTFConfig method
type AppTFConfigInput struct {
	// Deployment specification of the application
	Spec *deployment.Spec

	// Target environment
	Env *environment.Environment

	// If target environment has domain enabled, the exact path on the system
	// of the file containing the domain's TF state.
	DomainTFStateFile string

	// The exact path on the system of the file containing the environment's
	// TF state.
	EnvTFStateFile string
}

func (i *Infrastructure) AppTFConfig(input *AppTFConfigInput) (map[string]string, error) {
	res := map[string]string{
		"terraform.tf": i.tfCoreConfig(),
		"app.tf":       i.appTfConfig(input, appTfTpl),
	}
	if input.Env.DomainEnabled() {
		res["app_dns.tf"] = i.appTfConfig(input, appDnsTfTpl)
	}
	return res, nil
}

func (i *Infrastructure) appTfConfig(in *AppTFConfigInput, tpl string) string {
	var b strings.Builder

	t := template.Must(template.New("").Parse(tpl))
	data := map[string]interface{}{
		"env_name":              in.Env.Name,
		"app_name":              in.Spec.App.Name,
		"env_tfstate_file":      in.EnvTFStateFile,
		"target_group_resource": "",
	}
	if in.Env.DomainEnabled() {
		data["target_group_resource"] = "aws_alb_target_group.alb_to_ecs_service.arn"
		data["domain_tfstate_file"] = in.DomainTFStateFile
		data["domain_name"] = in.Env.Domain
	}

	t.Execute(&b, data)
	return b.String()
}

func (i *Infrastructure) CreateApplication(ctx context.Context, s *deployment.Spec, tf *tfexec.Terraform) error {
	if err := tf.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize terraform: %v", err)
	}
	if err := i.applyAppConfig(ctx, s, tf); err != nil {
		return fmt.Errorf("failed to apply terraform changes: %v", err)
	}
	return nil
}

func (i *Infrastructure) ModifyApplication(
	ctx context.Context, spec *deployment.Spec, tf *tfexec.Terraform,
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
