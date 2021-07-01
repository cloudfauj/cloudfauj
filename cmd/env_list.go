package cmd

import (
	"fmt"
	"github.com/cloudfauj/cloudfauj/api"
	"os"

	"github.com/spf13/cobra"
)

var envListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Environments",
	Long:  "This command returns a list of all Environments managed by Cloudfauj.",
	Run:   runEnvListCmd,
}

func runEnvListCmd(cmd *cobra.Command, args []string) {
	serverAddr, _ := cmd.Flags().GetString("server-addr")
	apiClient := api.NewClient(serverAddr)

	res, err := apiClient.ListEnvironments()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "an error occured while fetching list of environments: %v", err)
		return
	}
	for name := range res {
		fmt.Println(name)
	}
}
