package api

import (
	"fmt"
	"github.com/cloudfauj/cloudfauj/server"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
	"time"
)

const Version = server.ApiV1Prefix

// API represents a client that can interact with a Cloudfauj Server REST API
type API struct {
	HttpClient *http.Client
	WsDialer   *websocket.Dialer
	baseURL    *url.URL
}

// ServerEvent is a single event created by the cloudfauj server.
// It can describe the event itself or an error that occured due to it.
type ServerEvent struct {
	Message string
	Err     error
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
