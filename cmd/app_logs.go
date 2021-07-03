package cmd

import (
	"fmt"
	"github.com/cloudfauj/cloudfauj/api"
	"github.com/spf13/cobra"
)

var appLogsCmd = &cobra.Command{
	Use:   "logs [flags] APP",
	Args:  cobra.ExactArgs(1),
	Short: "Get Application logs",
	Long: `
    This command displays the logs produced by an Application.`,
	RunE:    runAppLogsCmd,
	Example: "cloudfauj app logs demo-server",
}

func runAppLogsCmd(cmd *cobra.Command, args []string) error {
	apiClient, err := api.NewClient(serverAddr)
	if err != nil {
		return err
	}

	fmt.Printf("Fetching logs for %s...\n", args[0])
	logs, err := apiClient.AppLogs(args[0])
	if err != nil {
		return err
	}
	for log := range logs {
		fmt.Println(log)
	}
	fmt.Println("Done")
	return nil
}
