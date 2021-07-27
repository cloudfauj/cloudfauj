package api

import (
	"encoding/json"
	"github.com/cloudfauj/cloudfauj/deployment"
)

// Deploy requests the server to deploy an application.
// It streams all the deployment logs.
func (a *API) Deploy(spec *deployment.Spec) (<-chan *ServerEvent, error) {
	m, _ := json.Marshal(spec)
	return a.makeWebsocketRequest(a.constructWsURL("/app/deploy"), m)
}
