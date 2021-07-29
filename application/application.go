package application

import (
	"errors"
	"strings"
)

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
			BindPort int32 `json:"bind_port" mapstructure:"bind_port"`
		} `json:"network"`
	} `json:"resources"`

	Visibility string `json:"visibility"`
}

func (a *Application) CheckIsValid() error {
	if len(strings.TrimSpace(a.Name)) == 0 {
		return errors.New("name cannot be empty")
	}
	if a.Type != TypeServer {
		return errors.New("only " + TypeServer + " type is supported")
	}
	if a.Visibility != VisibilityPublic {
		return errors.New("only " + VisibilityPublic + " visibility is supported")
	}
	return nil
}
