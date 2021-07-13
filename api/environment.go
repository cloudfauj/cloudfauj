package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (a *API) CreateEnvironment(config map[string]interface{}) (<-chan *ServerEvent, error) {
	return a.makeWebsocketRequest(a.constructWsURL("/environment/create"), config)
}

func (a *API) DestroyEnvironment(name string) (<-chan *ServerEvent, error) {
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
