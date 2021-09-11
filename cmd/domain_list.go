package cmd

import (
	"fmt"
	"github.com/cloudfauj/cloudfauj/api"
	"github.com/spf13/cobra"
)

var domainListCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all domains registered in Cloudfauj",
	Long: `
    This command returns a list of domains registered in Cloudfauj.

    These domains can be associated with environments to assign TLS-enabled
    URLs to the applications deployed in them.`,
	RunE:    runDomainListCmd,
	Example: "cloudfauj domain ls",
}

func runDomainListCmd(cmd *cobra.Command, args []string) error {
	apiClient, err := api.NewClient(serverAddr)
	if err != nil {
		return err
	}
	res, err := apiClient.ListDomains()
	if err != nil {
		return err
	}
	if len(res) == 0 {
		fmt.Println("No domains registered yet")
	}
	for _, name := range res {
		fmt.Println(name)
	}
	return nil
}
