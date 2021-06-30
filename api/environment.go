package api

type Environment struct {
	name string
}

func (a *API) CreateEnvironment() error {
	return nil
}

func (a *API) GetEnvironment(name string) (*Environment, error) {
	return nil, nil
}

func (a *API) ListEnvironments() ([]string, error) {
	return nil, nil
}
