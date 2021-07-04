package api

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
	"path"
	"time"
)

const Version = "v1"

// API represents a client that can interact with a Cloudfauj Server REST API
type API struct {
	HttpClient *http.Client
	WsDialer   *websocket.Dialer
	baseURL    *url.URL
}

// NewClient returns a new, initialized client to interact with a Cloudfauj Server.
func NewClient(serverAddr string) (*API, error) {
	// the baseURL we set must always contain at least scheme & hostname
	u, err := url.Parse(serverAddr)
	if err != nil {
		return nil, fmt.Errorf("invalid server url %s: %v", serverAddr, err)
	}

	return &API{
		HttpClient: &http.Client{Timeout: time.Minute},
		WsDialer:   websocket.DefaultDialer,
		baseURL:    u,
	}, nil
}

func (a *API) constructWsURL(p string) string {
	return a.constructURL("ws", p)
}

func (a *API) constructHttpURL(p string) string {
	return a.constructURL(a.baseURL.Scheme, p)
}

func (a *API) constructURL(s, p string) string {
	u := url.URL{Scheme: s, Host: a.baseURL.Host, Path: path.Join("/"+Version, p)}
	return u.String()
}
