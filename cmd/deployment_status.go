package cmd

import (
	"fmt"
	"github.com/cloudfauj/cloudfauj/api"
	"os"

	"github.com/spf13/cobra"
)

var deploymentStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get information about a Deployment",
	Long: `This command displays information about a deployment.
Among other things, it returns the status of the deployment.

Example:
    cloudfauj deployment status 123456`,
	Run: runDeploymentStatusCmd,
}

func runDeploymentStatusCmd(cmd *cobra.Command, args []string) {
	serverAddr, _ := cmd.Flags().GetString("server-addr")
	apiClient := api.NewClient(serverAddr)

	res, err := apiClient.Deployment(args[0])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "an error occured while fetching deployment status: %v", err)
		return
	}
	fmt.Println(res)
}
