package cmd

import (
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Launch a Cloudfauj Server",
	Long: `
    This command starts a Cloudfauj Server that carries out tasks such
    as Deployments when requested.
    The server takes care of provisioning and managing all infrastructure required
    by the application.`,
	RunE: runServerCmd,
}

func runServerCmd(cmd *cobra.Command, args []string) error {
	/*
		1. Load server config file
		2. Start server and block
	*/
	return nil
}
