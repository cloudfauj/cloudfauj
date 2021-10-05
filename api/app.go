package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudfauj/cloudfauj/deployment"
	"github.com/cloudfauj/cloudfauj/server"
	"net/http"
)

// Deploy requests the server to deploy an application.
// It streams all the deployment logs.
func (a *API) Deploy(spec *deployment.Spec) (<-chan *server.Event, error) {
	m, _ := json.Marshal(spec)
	return a.makeWebsocketRequest(a.constructWsURL("/app/deploy"), m)
}

func (a *API) DestroyApp(app, env string) error {
	u := a.constructHttpURL("/app/"+app, qp{"env": env})
	req, _ := http.NewRequest(http.MethodDelete, u, nil)

	res, err := a.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusNotFound:
		return errors.New("the target app or environment does not exist")
	case http.StatusOK:
		return nil
	}

	return fmt.Errorf("server returned %d: %v", res.StatusCode, err)
}
