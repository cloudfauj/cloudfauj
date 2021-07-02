package cmd

import (
	"github.com/spf13/cobra"
)

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Manage Applications",
	Long: `
    This command lets you manage and interact with Applications.`,
}
