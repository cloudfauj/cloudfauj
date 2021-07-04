package cmd

import (
	"fmt"
	"github.com/cloudfauj/cloudfauj/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var deployCmd = &cobra.Command{
	Use:   "deploy --env ENV [flags] ARTIFACT",
	Args:  cobra.ExactArgs(1),
	Short: "Deploy an Application",
	Long: `
    This command lets you deploy an Application to an environment.

    By default, it looks for the .cloudfauj.yml file in the current
    directory as the app configuration.

    As of today, the only type of application supported is a Docker
    container running a TCP server.
    The value of ARTIFACT must be the URI of a docker image residing
    in AWS ECR.`,
	RunE:    runDeployCmd,
	Example: "cloudfauj deploy --env staging 123456789012.dkr.ecr.us-east-1.amazonaws.com/demo-server:v1.0.0",
}

func init() {
	deployCmd.Flags().String("config", ".cloudfauj.yml", "Application configuration file")
	deployCmd.Flags().String("env", "", "The environment to deploy to")
	_ = deployCmd.MarkFlagRequired("env")
}

func runDeployCmd(cmd *cobra.Command, args []string) error {
	apiClient, err := api.NewClient(serverAddr)
	if err != nil {
		return err
	}

	configFile, _ := cmd.Flags().GetString("config")
	initConfig(configFile)
	viper.Set("artifact", args[0])

	targetEnv, _ := cmd.Flags().GetString("env")
	viper.Set("target_env", targetEnv)

	fmt.Printf("Deploying %s (%s) to %s\n", viper.GetString("name"), args[0], targetEnv)
	eventsCh, err := apiClient.Deploy(cmd.Context(), viper.AllSettings())
	if err != nil {
		return err
	}

	fmt.Println("Streaming deployment logs from server...")
	for e := range eventsCh {
		if e.Err != nil {
			return e.Err
		}
		fmt.Println(e.Message)
	}

	return nil
}
