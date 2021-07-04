package api

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
)

type DeploymentEvent struct {
	Message string
	Err     error
}

// Deploy requests the server to deploy an application.
// It streams all the deployment logs.
func (a *API) Deploy(ctx context.Context, appSpec map[string]interface{}) (<-chan *DeploymentEvent, error) {
	eventsCh := make(chan *DeploymentEvent)

	conn, _, err := a.WsDialer.Dial(a.constructWsURL("/app/deploy"), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to establish websocket connection with server: %v", err)
	}

	go func(conn *websocket.Conn, respCh chan<- *DeploymentEvent) {
		defer conn.Close()
		defer close(respCh)

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				// unless an error has occurred due to normal connection closure
				// from server, it needs to propagate.
				if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					respCh <- &DeploymentEvent{
						Err: fmt.Errorf("unexpected error during deployment: %v", err),
					}
				}
				break
			}
			respCh <- &DeploymentEvent{Message: string(msg)}
		}
	}(conn, eventsCh)

	if err = conn.WriteJSON(appSpec); err != nil {
		return nil, fmt.Errorf("failed to send data to server: %v", err)
	}
	return eventsCh, nil
}

func (a *API) AppLogs(app string) ([]string, error) {
	return nil, nil
}
