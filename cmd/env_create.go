package cmd

import (
	"fmt"
	"github.com/cloudfauj/cloudfauj/api"
	"os"

	"github.com/spf13/cobra"
)

// envCreateCmd represents the create command
var envCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create new Environment",
	Long: `This command lets you create a new environment.

An environment is a group of applications that's logically isolated from other
environments. It is created and managed by a CloudFauj admin.

At least 1 env must exist for an application to be deployed.`,
	Run: runEnvCreateCmd,
}

func init() {
	envCreateCmd.Flags().String(
		"config",
		".cloudfauj-env.yml",
		"Configuration file to create an environment from",
	)
}

func runEnvCreateCmd(cmd *cobra.Command, args []string) {
	serverAddr, _ := cmd.Flags().GetString("server-addr")
	apiClient := api.NewClient(serverAddr)

	configFile, _ := cmd.Flags().GetString("config")
	initConfig(configFile)

	if err := apiClient.CreateEnvironment(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "an error occured while creating environment: %v", err)
		return
	}
}
