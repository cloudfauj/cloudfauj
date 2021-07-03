package api

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
	"path"
	"time"
)

// API represents a client that can interact with a Cloudfauj Server REST API
type API struct {
	HttpClient *http.Client
	WsDialer   *websocket.Dialer
	baseURL    *url.URL
}

// NewClient returns a new, initialized client to interact with a Cloudfauj Server.
func NewClient(serverAddr string) (*API, error) {
	u, err := url.Parse(serverAddr)
	if err != nil {
		return nil, fmt.Errorf("invalid server url %s: %v", serverAddr, err)
	}
	u.Path = path.Join(u.Path, "/v1")

	return &API{
		HttpClient: &http.Client{Timeout: time.Minute},
		WsDialer:   websocket.DefaultDialer,
		baseURL:    u,
	}, nil
}

func (a *API) serverWebsocketURL(p string) string {
	u := url.URL{Scheme: "ws", Host: a.baseURL.Host, Path: path.Join(a.baseURL.Path, p)}
	return u.String()
}
