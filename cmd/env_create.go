package cmd

import (
	"fmt"
	"github.com/cloudfauj/cloudfauj/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var envCreateCmd = &cobra.Command{
	Use:   "create --config PATH",
	Short: "Create a new Environment",
	Long: `
    This command lets you create a new environment.

    You must provide a configuration to create the env from.
    The config defines the underlying infrastructure to provision to manage resources
    such as DNS, applications, etc.

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

	envName := viper.GetString("name")

	fmt.Printf("Creating environment %s\n", envName)
	if err := apiClient.CreateEnvironment(envName, viper.AllSettings()); err != nil {
		return err
	}

	fmt.Println("Done")
	return nil
}
