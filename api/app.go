package api

import (
	"encoding/json"
	"fmt"
	"github.com/cloudfauj/cloudfauj/deployment"
	"net/http"
)

// Deploy requests the server to deploy an application.
// It streams all the deployment logs.
func (a *API) Deploy(spec *deployment.Spec) (<-chan *ServerEvent, error) {
	m, _ := json.Marshal(spec)
	return a.makeWebsocketRequest(a.constructWsURL("/app/deploy"), m)
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
