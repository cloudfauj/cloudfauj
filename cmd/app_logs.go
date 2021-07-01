package cmd

import (
	"fmt"
	"github.com/cloudfauj/cloudfauj/api"
	"github.com/spf13/cobra"
	"os"
)

var appLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Get Application logs",
	Long: `This command displays the logs produced by an Application.
The example below fetches logs for an app called "demo-server":
    cloudfauj app logs demo-server`,
	Run: runAppLogsCmd,
}

func runAppLogsCmd(cmd *cobra.Command, args []string) {
	serverAddr, _ := cmd.Flags().GetString("server-addr")
	apiClient := api.NewClient(serverAddr)

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
