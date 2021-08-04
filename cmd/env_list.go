package cmd

import (
	"fmt"
	"github.com/cloudfauj/cloudfauj/api"
	"github.com/spf13/cobra"
)

var envListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Environments",
	Long: `
    This command returns a list of all Environments managed by Cloudfauj.`,
	RunE: runEnvListCmd,
}

func runEnvListCmd(cmd *cobra.Command, args []string) error {
	apiClient, err := api.NewClient(serverAddr)
	if err != nil {
		return err
	}
	res, err := apiClient.ListEnvironments()
	if err != nil {
		return err
	}
	if len(res) == 0 {
		fmt.Println("No environments created yet")
	}
	for _, name := range res {
		fmt.Println(name)
	}
	return nil
}
