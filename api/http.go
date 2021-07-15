package api

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/url"
	"path"
)

// queryParams holds URL query parameters in a structured way
type queryParams map[string]string

func (a *API) constructWsURL(p string) string {
	return a.constructURL("ws", p, nil)
}

func (a *API) constructHttpURL(p string, q queryParams) string {
	return a.constructURL(a.baseURL.Scheme, p, q)
}

func (a *API) constructURL(s, p string, q queryParams) string {
	u := url.URL{Scheme: s, Host: a.baseURL.Host, Path: path.Join("/"+Version, p)}
	// set query parameters while preserving any previous ones
	if q != nil {
		original := u.Query()
		for p, v := range q {
			original.Set(p, v)
		}
		u.RawQuery = original.Encode()
	}
	return u.String()
}

func (a *API) makeWebsocketRequest(u string, message []byte) (<-chan *ServerEvent, error) {
	eventsCh := make(chan *ServerEvent)

	conn, _, err := a.WsDialer.Dial(u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to establish websocket connection with server: %v", err)
	}

	go func(conn *websocket.Conn, respCh chan<- *ServerEvent) {
		defer conn.Close()
		defer close(respCh)

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				// unless an error has occurred due to normal connection closure
				// from server, it needs to propagate.
				if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
					respCh <- &ServerEvent{
						Err: fmt.Errorf("server error: %v", err),
					}
				}
				break
			}
			respCh <- &ServerEvent{Message: string(msg)}
		}
	}(conn, eventsCh)

	if message != nil {
		if err = conn.WriteMessage(websocket.BinaryMessage, message); err != nil {
			return nil, fmt.Errorf("failed to send data to server: %v", err)
		}
	}
	return eventsCh, nil
}
