package cmd

import (
	"github.com/spf13/cobra"
)

var deploymentCmd = &cobra.Command{
	Use:   "deployment",
	Short: "Manage Deployments",
	Long:  "This command helps you manage and interact with Deployments.",
}
