package cmd

import (
	"fmt"
	"github.com/cloudfauj/cloudfauj/api"
	"github.com/spf13/cobra"
)

var deploymentListCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all Deployments",
	Long: `
    This command displays a list of all running deployments in Cloudfauj`,
	RunE: runDeploymentListCmd,
}

func runDeploymentListCmd(cmd *cobra.Command, args []string) error {
	apiClient, err := api.NewClient(serverAddr)
	if err != nil {
		return err
	}
	res, err := apiClient.ListDeployments()
	if err != nil {
		return err
	}
	if len(res) == 0 {
		fmt.Println("No deployments running at this time")
	}
	for _, d := range res {
		desc := `ID: %s
    App:        %s
    Target Env: %s
    Status:     %s

`
		fmt.Printf(desc, d.Id, d.App, d.Environment, d.Status)
	}
	return nil
}
