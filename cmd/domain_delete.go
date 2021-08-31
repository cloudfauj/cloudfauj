package cmd

import (
	"fmt"
	"github.com/cloudfauj/cloudfauj/api"
	"github.com/spf13/cobra"
)

var domainDeleteCmd = &cobra.Command{
	Use:   "delete [flags] DOMAIN",
	Short: "De-register a domain from Cloudfauj",
	Long: `
    This command lets you delete a domain previously added to Cloudfauj.

    All infrastructure for managing this domain, including the TLS certificates,
    is destroyed. Any applications with a URL on this domain will no longer be
    reachable. 

    Be very careful with this command!`,
	Args:    cobra.ExactArgs(1),
	RunE:    runDomainDeleteCmd,
	Example: "cloudfauj domain delete example.com",
}

func runDomainDeleteCmd(cmd *cobra.Command, args []string) error {
	apiClient, err := api.NewClient(serverAddr)
	if err != nil {
		return err
	}
	eventsCh, err := apiClient.DeleteDomain(args[0])
	if err != nil {
		return err
	}
	for e := range eventsCh {
		if e.Err != nil {
			return e.Err
		}
		fmt.Println(e.Msg)
	}
	return nil
}
