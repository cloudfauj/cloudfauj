package deployment

type Deployment struct {
	Id          string `json:"id"`
	App         string `json:"app"`
	Environment string `json:"environment"`
	Status      string `json:"status"`
}
