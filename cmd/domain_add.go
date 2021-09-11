package cmd

import (
	"fmt"
	"github.com/cloudfauj/cloudfauj/api"
	"github.com/spf13/cobra"
)

var domainAddCmd = &cobra.Command{
	Use:   "add [flags] DOMAIN",
	Short: "Register a domain with Cloudfauj",
	Long: `
    This command lets you add a domain to Cloudfauj.

    Adding a domain is necessary before you can start using it to assign URLs to apps.
    Upon adding, Cloudfauj creates some AWS infrastructure like ACM Certificates and
    Route53 Hosted Zone to manage URLs.

    This command outputs NS records of the hosted zone that need to be configured for
    your domain in your domain provider's dashboard.`,
	Args:    cobra.ExactArgs(1),
	RunE:    runDomainAddCmd,
	Example: "cloudfauj domain add example.com",
}

func runDomainAddCmd(cmd *cobra.Command, args []string) error {
	apiClient, err := api.NewClient(serverAddr)
	if err != nil {
		return err
	}
	eventsCh, err := apiClient.AddDomain(args[0])
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
