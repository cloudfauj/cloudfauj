package api

type Deployment struct{}

func (a *API) DeploymentStatus(id string) (*Deployment, error) {
	return nil, nil
}

func (a *API) DeploymentLogs(id string) ([]string, error) {
	return nil, nil
}
