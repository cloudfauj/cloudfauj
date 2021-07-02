package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var envListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Environments",
	Long: `
    This command returns a list of all Environments managed by Cloudfauj.`,
	Run: runEnvListCmd,
}

func runEnvListCmd(cmd *cobra.Command, args []string) {
	apiClient := createApiClient()

	res, err := apiClient.ListEnvironments()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "an error occured while fetching list of environments: %v", err)
		return
	}
	for name := range res {
		fmt.Println(name)
	}
}
