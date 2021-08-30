package cmd

import "github.com/spf13/cobra"

var domainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Manage Custom Domains",
	Long: `
    This command lets you manage how Cloudfauj uses your custom Domain(s)
    to assign TLS-enabled URLs to applications.

    You can add domains to Cloudfauj, then associate them with environments.
    Every public application deployed to a domain-enabled environment receives
    a TLS-enabled URL which you can use to access it.

    Domains are an optional feature.`,
}
