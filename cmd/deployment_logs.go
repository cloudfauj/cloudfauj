package cmd

import (
	"fmt"
	"github.com/cloudfauj/cloudfauj/api"
	"os"

	"github.com/spf13/cobra"
)

// deploymentLogsCmd represents the logs command
var deploymentLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: runDeploymentLogsCmd,
}

func runDeploymentLogsCmd(cmd *cobra.Command, args []string) {
	serverAddr, _ := cmd.Flags().GetString("server-addr")
	apiClient := api.NewClient(serverAddr)

	res, err := apiClient.DeploymentLogs(args[0])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "an error occured while deploying: %v", err)
		return
	}
	fmt.Println(res)
}
