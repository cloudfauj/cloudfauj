package deployment

const (
	StatusCreated   = "created"
	StatusOngoing   = "ongoing"
	StatusSucceeded = "succeeded"
	StatusFailed    = "failed"
)

type Deployment struct {
	Id          string `json:"id"`
	App         string `json:"app"`
	Environment string `json:"environment"`
	Status      string `json:"status"`
}
