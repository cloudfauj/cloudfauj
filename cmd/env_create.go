package cmd

import (
	"fmt"
	"github.com/cloudfauj/cloudfauj/api"
	"github.com/cloudfauj/cloudfauj/environment"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var envCreateCmd = &cobra.Command{
	Use:   "create --config PATH",
	Short: "Create a new Environment",
	Long: `
    This command lets you create a new environment.

    You must provide a configuration to create the environment from.
    The config defines the underlying infrastructure to provision to manage different
    types of resources such as container orchestrator.

    At least 1 env must exist for an application to be deployed.`,
	RunE:    runEnvCreateCmd,
	Example: "cloudfauj env create --config ./cloudfauj-env.yml",
}

func init() {
	envCreateCmd.Flags().String("config", "", "Configuration file to create an environment from")
	_ = envCreateCmd.MarkFlagRequired("config")
}

func runEnvCreateCmd(cmd *cobra.Command, args []string) error {
	apiClient, err := api.NewClient(serverAddr)
	if err != nil {
		return err
	}
	configFile, _ := cmd.Flags().GetString("config")
	initConfig(configFile)

	var env environment.Environment
	_ = viper.Unmarshal(&env)

	fmt.Printf("Requesting creation of %s\n\n", env.Name)
	eventsCh, err := apiClient.CreateEnvironment(&env)
	if err != nil {
		return err
	}
	for e := range eventsCh {
		if e.Err != nil {
			return e.Err
		}
		fmt.Println(e.Message)
	}
	return nil
}
