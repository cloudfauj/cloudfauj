package deployment

import (
	"github.com/sirupsen/logrus"
)

const (
	StatusRunning   = "running"
	StatusSucceeded = "succeeded"
	StatusFailed    = "failed"
)

type Deployment struct {
	Id          string `json:"id"`
	App         string `json:"app"`
	Environment string `json:"environment"`
	Status      string `json:"status"`
	log         *logrus.Logger
}

func New(s *Spec, l *logrus.Logger) *Deployment {
	return &Deployment{
		App:         s.App.Name,
		Environment: s.TargetEnv,
		Status:      StatusRunning,
		log:         l,
	}
}

func (d *Deployment) Log(msg string) {
	d.log.Info(msg)
}

func (d *Deployment) Fail(err error) {
	d.Status = StatusFailed
	d.log.Errorf("Deployment Failed: %v", err)
}

func (d *Deployment) Succeed() {
	d.Status = StatusSucceeded
	d.log.Info("Deployment successful")
}
