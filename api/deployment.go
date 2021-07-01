package api

type Deployment struct {
	id          string
	app         string
	environment string
	status      string
}

func (a *API) Deployment(id string) (*Deployment, error) {
	return nil, nil
}

func (a *API) DeploymentLogs(id string) ([]string, error) {
	return nil, nil
}

func (a *API) ListDeployments() ([]*Deployment, error) {
	return nil, nil
}
