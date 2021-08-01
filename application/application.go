package application

import (
	"errors"
	"strings"
)

const TypeServer = "server"

const VisibilityPublic = "public"

type Application struct {
	Name        string       `json:"name"`
	Type        string       `json:"type"`
	Visibility  string       `json:"visibility"`
	HealthCheck *HealthCheck `json:"healthcheck"`
	Resources   *Resources   `json:"resources"`
}

type HealthCheck struct {
	Path string `json:"path"`
}

type Resources struct {
	Cpu     int      `json:"cpu"`
	Memory  int      `json:"memory"`
	Network *Network `json:"network"`
}

type Network struct {
	BindPort int32 `json:"bind_port" mapstructure:"bind_port"`
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
