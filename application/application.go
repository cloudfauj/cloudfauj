package application

const TypeServer = "server"

const VisibilityPublic = "public"

type Application struct {
	Name string `json:"name"`
	Type string `json:"type"`

	HealthCheck struct {
		Path string `json:"path"`
	} `json:"healthcheck"`

	Resources struct {
		Cpu    int `json:"cpu"`
		Memory int `json:"memory"`

		Network struct {
			BindPort int `json:"bind_port" mapstructure:"bind_port"`
		} `json:"network"`
	} `json:"resources"`

	Visibility string `json:"visibility"`
}

func (a *Application) CheckIsValid() error {
	return nil
}
