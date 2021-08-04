package cmd

import (
	"fmt"
	"github.com/cloudfauj/cloudfauj/api"
	"github.com/cloudfauj/cloudfauj/application"
	"github.com/cloudfauj/cloudfauj/deployment"
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
    container running a HTTP/TCP server in AWS ECS-Fargate.
    The value of ARTIFACT must be the URI of a docker image residing
    in AWS ECR.`,
	RunE:    runDeployCmd,
	Example: "cloudfauj deploy --env staging 123456789012.dkr.ecr.us-east-1.amazonaws.com/demo-server:latest",
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

	var app application.Application
	_ = viper.Unmarshal(&app)

	env, _ := cmd.Flags().GetString("env")
	spec := &deployment.Spec{App: &app, TargetEnv: env, Artifact: args[0]}

	fmt.Printf("Deploying %s artifact to %s\n\n", app.Name, spec.TargetEnv)
	eventsCh, err := apiClient.Deploy(spec)
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
