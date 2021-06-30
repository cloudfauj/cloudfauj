package api

// API represents a client that can interact with a Cloudfauj server REST API
type API struct {
	server string
}

// NewClient returns a new, initialized client to interact with a Cloudfauj Server.
func NewClient(serverAddr string) *API {
	return &API{server: serverAddr}
}
