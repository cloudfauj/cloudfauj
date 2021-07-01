package api

// API represents a client that can interact with a Cloudfauj serverAddr REST API
type API struct {
	serverAddr string
}

// NewClient returns a new, initialized client to interact with a Cloudfauj Server.
func NewClient(serverAddr string) *API {
	return &API{serverAddr}
}
