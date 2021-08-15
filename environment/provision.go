package environment

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-exec/tfexec"
	"os"
)

const EnvConfig = `module "%s" {
  source              = "%s"
  env_name            = "%s"
  main_vpc_cidr_block = "%s"
}

output "%s_ecs_cluster_arn" {
  value = module.%s.compute_ecs_cluster_arn
}`

const EnvModuleSource = "github.com/cloudfauj/terraform-template.git//env?ref=90fceb4"

func (e *Environment) Provision(ctx context.Context, tfFile *os.File, eventsCh chan<- Event) {
	defer close(eventsCh)

	cidr, err := e.Infra.NextAvailableCIDR(ctx)
	if err != nil {
		eventsCh <- Event{Err: fmt.Errorf("failed to compute VPC CIDR: %v", err)}
		return
	}
	envConf := fmt.Sprintf(EnvConfig, e.Name, EnvModuleSource, e.Name, cidr, e.Name, e.Name)
	if _, err := tfFile.Write([]byte(envConf)); err != nil {
		eventsCh <- Event{Err: fmt.Errorf("failed to write Terraform configuration to file: %v", err)}
		return
	}
	if err := e.Infra.Tf.Init(ctx); err != nil {
		eventsCh <- Event{Err: fmt.Errorf("failed to initialize terraform: %v", err)}
		return
	}
	eventsCh <- Event{Msg: "Applying Terraform configuration"}
	if err := e.Infra.Tf.Apply(ctx, tfexec.Target("module."+e.Name)); err != nil {
		eventsCh <- Event{Err: fmt.Errorf("failed to apply TF config: %v", err)}
		return
	}
}
