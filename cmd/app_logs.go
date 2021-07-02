package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var appLogsCmd = &cobra.Command{
	Use:   "logs [flags] APP",
	Args:  cobra.ExactArgs(1),
	Short: "Get Application logs",
	Long: `
    This command displays the logs produced by an Application.`,
	Run:     runAppLogsCmd,
	Example: "cloudfauj app logs demo-server",
}

func runAppLogsCmd(cmd *cobra.Command, args []string) {
	apiClient := createApiClient()

	fmt.Printf("Fetching logs for %s...\n", args[0])
	logs, err := apiClient.Logs(args[0])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "an error occured while fetching logs: %v", err)
		return
	}
	for log := range logs {
		fmt.Println(log)
	}
	fmt.Println("Done")
}
