package deployment

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
}

func New(s *Spec) *Deployment {
	return &Deployment{
		App:         s.App.Name,
		Environment: s.TargetEnv,
		Status:      StatusRunning,
	}
}

func (d *Deployment) AppendLog(m string) error {
	// add log to logfile
	return nil
}

func (d *Deployment) Fail(e error) error {
	d.Status = StatusFailed
	// write failure log to logfile
	return nil
}
