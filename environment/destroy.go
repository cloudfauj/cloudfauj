package environment

import "context"

func (e *Environment) Destroy(ctx context.Context, eventsCh chan<- Event) {
	defer close(eventsCh)
}
