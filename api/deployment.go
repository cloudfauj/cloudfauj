package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudfauj/cloudfauj/deployment"
	"net/http"
)

func (a *API) Deployment(id string) (*deployment.Deployment, error) {
	var result deployment.Deployment

	res, err := a.HttpClient.Get(a.constructHttpURL("/deployment/"+id, nil))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil, errors.New("deployment does not exist")
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned %d: %v", res.StatusCode, err)
	}
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode server response: %v", err)
	}
	return &result, nil
}

func (a *API) DeploymentLogs(id string) ([]string, error) {
	var result []string

	res, err := a.HttpClient.Get(a.constructHttpURL("/deployment/"+id+"/logs", nil))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil, errors.New("logs for deployment do not exist")
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned %d: %v", res.StatusCode, err)
	}
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode server response: %v", err)
	}
	return result, nil
}

func (a *API) ListDeployments() ([]*deployment.Deployment, error) {
	var result []*deployment.Deployment

	res, err := a.HttpClient.Get(a.constructHttpURL("/deployments", nil))
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
