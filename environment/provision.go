package environment

import (
	"context"
	"fmt"
	"os"
)

func (e *Environment) Provision(ctx context.Context, tfDir string, tfFile *os.File, eventsCh chan<- Event) {
	defer close(eventsCh)

	tf, err := e.Infra.Tf(e.Name)
	if err != nil {
		eventsCh <- Event{Err: err}
		return
	}
	cidr, err := e.Infra.NextAvailableCIDR(ctx)
	if err != nil {
		eventsCh <- Event{Err: fmt.Errorf("failed to compute VPC CIDR: %v", err)}
		return
	}
	envConf := e.Infra.EnvTfConfig(e.Name, cidr)
	if _, err := tfFile.Write([]byte(envConf)); err != nil {
		eventsCh <- Event{Err: fmt.Errorf("failed to write Terraform configuration to file: %v", err)}
		return
	}
	if err := tf.Init(ctx); err != nil {
		eventsCh <- Event{Err: fmt.Errorf("failed to initialize terraform: %v", err)}
		return
	}
	eventsCh <- Event{Msg: "Applying Terraform configuration"}
	if err := tf.Apply(ctx); err != nil {
		eventsCh <- Event{Err: fmt.Errorf("failed to apply TF config: %v", err)}
		return
	}
}
