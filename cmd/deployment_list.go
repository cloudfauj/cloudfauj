package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var deploymentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Deployments",
	Long: `
    This command displays a list of all Deployments active in Cloudfauj`,
	RunE: runDeploymentListCmd,
}

func runDeploymentListCmd(cmd *cobra.Command, args []string) error {
	apiClient := createApiClient()

	res, err := apiClient.ListDeployments()
	if err != nil {
		return err
	}
	for d := range res {
		fmt.Println(d)
	}
	return nil
}
