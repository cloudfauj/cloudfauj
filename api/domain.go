package api

import (
	"encoding/json"
	"fmt"
	"github.com/cloudfauj/cloudfauj/domain"
	"github.com/cloudfauj/cloudfauj/server"
	"net/http"
)

func (a *API) AddDomain(d *domain.Domain) (<-chan *server.Event, error) {
	m, _ := json.Marshal(d)
	return a.makeWebsocketRequest(a.constructWsURL("/domain/add"), m)
}

func (a *API) DeleteDomain(name string) (<-chan *server.Event, error) {
	return a.makeWebsocketRequest(a.constructWsURL("/domain/"+name+"/delete"), nil)
}

func (a *API) ListDomains() ([]string, error) {
	var result []string

	res, err := a.HttpClient.Get(a.constructHttpURL("/domains", nil))
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

func (a *API) TFPlanDomain(name string) (<-chan *server.Event, error) {
	return a.makeWebsocketRequest(a.constructWsURL("/domain/"+name+"/plan"), nil)
}
