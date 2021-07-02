package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var deploymentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Deployments",
	Long: `
    This command displays a list of all Deployments active in Cloudfauj`,
	Run: runDeploymentListCmd,
}

func runDeploymentListCmd(cmd *cobra.Command, args []string) {
	apiClient := createApiClient()

	res, err := apiClient.ListDeployments()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "an error occured while fetching deployments: %v", err)
		return
	}
	for d := range res {
		fmt.Println(d)
	}
}
