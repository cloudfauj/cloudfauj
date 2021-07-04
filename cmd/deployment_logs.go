package cmd

import (
	"fmt"
	"github.com/cloudfauj/cloudfauj/api"
	"github.com/spf13/cobra"
)

// deploymentLogsCmd represents the logs command
var deploymentLogsCmd = &cobra.Command{
	Use:   "logs [flags] ID",
	Short: "Fetch deployment logs",
	Long: `
    This command displays logs of a deployment.
    You must specify a deployment ID to fetch logs of.`,
	Args:    cobra.ExactArgs(1),
	RunE:    runDeploymentLogsCmd,
	Example: "cloudfauj deployment logs 123456",
}

func runDeploymentLogsCmd(cmd *cobra.Command, args []string) error {
	apiClient, err := api.NewClient(serverAddr)
	if err != nil {
		return err
	}

	logs, err := apiClient.DeploymentLogs(args[0])
	if err != nil {
		return err
	}
	for _, log := range logs {
		fmt.Println(log)
	}
	return nil
}
