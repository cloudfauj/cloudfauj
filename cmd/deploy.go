package cmd

import (
	"fmt"
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
    in AWS ECR. 

    This command returns a Deployment ID that you can use to fetch the status
    and logs of the deployment.`,
	RunE:    runDeployCmd,
	Example: "cloudfauj deploy --env staging 123456789012.dkr.ecr.us-east-1.amazonaws.com/demo-server:v1.0.0",
}

func init() {
	deployCmd.Flags().String("config", ".cloudfauj.yml", "Application configuration file")
	deployCmd.Flags().String("env", "", "The environment to deploy to")
	_ = deployCmd.MarkFlagRequired("env")
}

func runDeployCmd(cmd *cobra.Command, args []string) error {
	apiClient := createApiClient()
	configFile, _ := cmd.Flags().GetString("config")
	initConfig(configFile)

	targetEnv, _ := cmd.Flags().GetString("env")
	appName := viper.GetString("name")

	fmt.Printf("Deploying %s (%s) to %s\n", appName, args[0], targetEnv)
	res, err := apiClient.Deploy(args[0], viper.AllSettings())
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}
