package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Deployment struct {
	Id          string `json:"id"`
	App         string `json:"app"`
	Environment string `json:"environment"`
	Status      string `json:"status"`
}

func (a *API) Deployment(id string) (*Deployment, error) {
	return nil, nil
}

func (a *API) DeploymentLogs(id string) ([]string, error) {
	var result []string

	res, err := a.HttpClient.Get(a.constructHttpURL("/deployment/" + id + "/logs"))
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

func (a *API) ListDeployments() ([]*Deployment, error) {
	var result []*Deployment

	res, err := a.HttpClient.Get(a.constructHttpURL("/deployments"))
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
