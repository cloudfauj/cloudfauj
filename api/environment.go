package api

import (
	"encoding/json"
	"fmt"
	"github.com/cloudfauj/cloudfauj/environment"
	"github.com/cloudfauj/cloudfauj/server"
	"net/http"
)

func (a *API) CreateEnvironment(env *environment.Environment) (<-chan *server.Event, error) {
	m, _ := json.Marshal(env)
	return a.makeWebsocketRequest(a.constructWsURL("/environment/create"), m)
}

func (a *API) DestroyEnvironment(name string) (<-chan *server.Event, error) {
	return a.makeWebsocketRequest(a.constructWsURL("/environment/"+name+"/destroy"), nil)
}

func (a *API) ListEnvironments() ([]string, error) {
	var result []string

	res, err := a.HttpClient.Get(a.constructHttpURL("/environments", nil))
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

func (a *API) TFPlanEnv(name string) (<-chan *server.Event, error) {
	return a.makeWebsocketRequest(a.constructWsURL("/environment/"+name+"/plan"), nil)
}

func (a *API) TFApplyEnv(name string) (<-chan *server.Event, error) {
	return a.makeWebsocketRequest(a.constructWsURL("/environment/"+name+"/apply"), nil)
}
