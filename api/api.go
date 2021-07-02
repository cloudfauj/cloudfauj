package api

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"
)

// API represents a client that can interact with a Cloudfauj Server REST API
type API struct {
	HttpClient *http.Client
	baseURL    string
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
		baseURL:    u.String(),
	}, nil
}
