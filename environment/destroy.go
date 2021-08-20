package environment

import (
	"context"
	"fmt"
)

func (e *Environment) Destroy(ctx context.Context, eventsCh chan<- Event) {
	defer close(eventsCh)

	tf, err := e.Infra.Tf(e.Name)
	if err != nil {
		eventsCh <- Event{Err: err}
		return
	}
	eventsCh <- Event{Msg: "Destroying Terraform infrastructure"}
	if err := tf.Destroy(ctx); err != nil {
		eventsCh <- Event{Err: fmt.Errorf("failed to destroy: %v", err)}
		return
	}
}
