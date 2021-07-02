package cmd

import (
	"fmt"
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
	apiClient := createApiClient()

	res, err := apiClient.ListEnvironments()
	if err != nil {
		return err
	}
	for name := range res {
		fmt.Println(name)
	}
	return nil
}
