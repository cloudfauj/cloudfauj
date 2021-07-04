package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

// Deploy requests the server to deploy an application.
// It streams all the deployment logs.
func (a *API) Deploy(appSpec map[string]interface{}) (<-chan *DeploymentEvent, error) {
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

func (a *API) AppLogs(app, env string) ([]string, error) {
	var result []string

	u := a.constructHttpURL("/app/"+app+"/logs", queryParams{"env": env})
	res, err := a.HttpClient.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned %d: %v", res.StatusCode, err)
	}
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode server response: %v", err)
	}
	return result, nil
}
