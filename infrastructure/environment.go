package infrastructure

import (
	"context"
	"fmt"
	"github.com/cloudfauj/cloudfauj/environment"
	"github.com/hashicorp/terraform-exec/tfexec"
	"strings"
	"text/template"
)

func (i *Infrastructure) EnvTFConfig(
	ctx context.Context, e *environment.Environment, dsf string,
) (map[string]string, error) {
	cidr, err := i.nextAvailableCIDR(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to compute VPC CIDR: %v", err)
	}
	res := map[string]string{
		"terraform.tf":    i.tfCoreConfig(),
		"network.tf":      i.envTfConfig(envNetworkTfTpl, e.Name, cidr),
		"orchestrator.tf": fmt.Sprintf(envOrchestratorTfTpl, e.Name),
	}
	if e.DomainEnabled() {
		res["domain.tf"] = fmt.Sprintf(envDomainStateTfTpl, dsf)
		res["load_balancer.tf"] = i.envTfConfig(envAlbTfTpl, e.Name, cidr)
	}
	return res, nil
}

func (i *Infrastructure) CreateEnvironment(ctx context.Context, tf *tfexec.Terraform) error {
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

func (i *Infrastructure) envTfConfig(tpl, env, cidr string) string {
	var b strings.Builder
	t := template.Must(template.New("").Parse(tpl))
	data := map[string]interface{}{"env_name": env, "vpc_cidr": cidr}

	t.Execute(&b, data)
	return b.String()
}
