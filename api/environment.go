package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (a *API) CreateEnvironment(name string, config map[string]interface{}) error {
	return nil
}

func (a *API) DestroyEnvironment(name string) error {
	return nil
}

func (a *API) ListEnvironments() ([]string, error) {
	var result []string

	res, err := a.HttpClient.Get(a.constructHttpURL("/environments"))
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
