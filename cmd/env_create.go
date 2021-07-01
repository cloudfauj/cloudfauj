package cmd

import (
	"fmt"
	"github.com/cloudfauj/cloudfauj/api"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

// envCreateCmd represents the create command
var envCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Environment",
	Long: `This command lets you create a new environment.

You must provide a configuration to create the env from.
This defines the underlying infrastructure to provision to manage resources
such as DNS, applications, etc.

At least 1 env must exist for an application to be deployed.`,
	Run: runEnvCreateCmd,
}

func init() {
	envCreateCmd.Flags().String("config", "", "Configuration file to create an environment from")
	_ = envCreateCmd.MarkFlagRequired("config")
}

func runEnvCreateCmd(cmd *cobra.Command, args []string) {
	serverAddr, _ := cmd.Flags().GetString("server-addr")
	apiClient := api.NewClient(serverAddr)

	configFile, _ := cmd.Flags().GetString("config")
	initConfig(configFile)

	envName := viper.GetString("name")

	fmt.Printf("Creating environment %s\n", envName)
	if err := apiClient.CreateEnvironment(envName, viper.AllSettings()); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "an error occured while creating the environment: %v", err)
		return
	}
	fmt.Println("Done")
}
