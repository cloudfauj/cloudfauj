package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (a *API) CreateEnvironment(name string, config map[string]interface{}) error {
	u := a.constructHttpURL("/environment/"+name, nil)
	body, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to create JSON payload from config: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, u, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	res, err := a.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned %d: %v", res.StatusCode, err)
	}
	return nil
}

func (a *API) DestroyEnvironment(name string) error {
	u := a.constructHttpURL("/environment/"+name, nil)
	req, err := http.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return err
	}

	res, err := a.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned %d: %v", res.StatusCode, err)
	}
	return nil
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
