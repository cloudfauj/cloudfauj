package cmd

import (
	"fmt"
	"github.com/cloudfauj/cloudfauj/api"
	"github.com/spf13/cobra"
)

var appLogsCmd = &cobra.Command{
	Use:   "logs --env ENV [flags] APP",
	Args:  cobra.ExactArgs(1),
	Short: "Get Application logs",
	Long: `
    This command displays the logs produced by the latest deployment of an
    Application in a specified environment.`,
	RunE:    runAppLogsCmd,
	Example: "cloudfauj app logs --env staging demo-server",
}

func init() {
	appLogsCmd.Flags().String("env", "", "The environment to fetch the app's logs from")
	_ = appLogsCmd.MarkFlagRequired("env")
}

func runAppLogsCmd(cmd *cobra.Command, args []string) error {
	env, _ := cmd.Flags().GetString("env")
	apiClient, err := api.NewClient(serverAddr)
	if err != nil {
		return err
	}

	fmt.Printf("Fetching logs for %s from %s...\n", args[0], env)
	logs, err := apiClient.AppLogs(args[0], env)
	if err != nil {
		return err
	}
	for _, log := range logs {
		fmt.Println(log)
	}
	return nil
}
