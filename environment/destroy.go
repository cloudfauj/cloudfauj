package environment

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-exec/tfexec"
)

func (e *Environment) Destroy(ctx context.Context, eventsCh chan<- Event) {
	defer close(eventsCh)

	eventsCh <- Event{Msg: "Destroying Terraform infrastructure"}
	if err := e.Infra.Tf.Destroy(ctx, tfexec.Target("module."+e.Name)); err != nil {
		eventsCh <- Event{Err: fmt.Errorf("failed to destroy: %v", err)}
		return
	}
}
