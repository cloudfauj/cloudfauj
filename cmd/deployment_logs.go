package cmd

import (
	"fmt"
	"github.com/cloudfauj/cloudfauj/api"
	"os"

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
	Run:     runDeploymentLogsCmd,
	Example: "cloudfauj deployment logs 123456",
}

func runDeploymentLogsCmd(cmd *cobra.Command, args []string) {
	serverAddr, _ := cmd.Flags().GetString("server-addr")
	apiClient := api.NewClient(serverAddr)

	logs, err := apiClient.DeploymentLogs(args[0])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "an error occured while fetching deployment logs: %v", err)
		return
	}
	for log := range logs {
		fmt.Println(log)
	}
	fmt.Println("Done")
}
