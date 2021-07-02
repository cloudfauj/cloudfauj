package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy --env ENV",
	Short: "Deploy an Application",
	Long: `
    This command lets you deploy an Application to an environment.

    By default, it looks for the .cloudfauj.yml file in the current
    directory as the app configuration.

    It returns a Deployment ID that you can use to fetch the status & logs
    of the deployment.`,
	Run:     runDeployCmd,
	Example: "cloudfauj deploy --env staging",
}

func init() {
	deployCmd.Flags().String("config", ".cloudfauj.yml", "Application configuration file")
	deployCmd.Flags().String("env", "", "The environment to deploy to")
	_ = deployCmd.MarkFlagRequired("env")
}

func runDeployCmd(cmd *cobra.Command, args []string) {
	apiClient := createApiClient()
	configFile, _ := cmd.Flags().GetString("config")
	initConfig(configFile)

	targetEnv, _ := cmd.Flags().GetString("env")
	appName := viper.GetString("name")

	fmt.Printf("Deploying application %s to %s\n", appName, targetEnv)
	res, err := apiClient.Deploy(appName, viper.AllSettings())
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "an error occured while deploying: %v", err)
		return
	}
	fmt.Println(res)
}
