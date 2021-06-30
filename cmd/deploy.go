package cmd

import (
	"fmt"
	"github.com/cloudfauj/cloudfauj/api"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Run a deployment",
	Long:  "This command asks the cloudfauj server to run a new deployment.",
	Run:   runDeployCmd,
}

func init() {
	deployCmd.Flags().String("config", ".cloudfauj.yml", "Project configuration file")
	deployCmd.Flags().String("env", "", "The environment to deploy to")
	_ = deployCmd.MarkFlagRequired("env")
}

func runDeployCmd(cmd *cobra.Command, args []string) {
	serverAddr, _ := cmd.Flags().GetString("server-addr")
	apiClient := api.NewClient(serverAddr)

	configFile, _ := cmd.Flags().GetString("config")
	initConfig(configFile)

	targetEnv, _ := cmd.Flags().GetString("env")
	project := viper.GetString("project")
	apps := viper.GetStringMap("applications")

	fmt.Printf("Deploying project %s to %s", project, targetEnv)
	res, err := apiClient.Deploy(project, apps)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "an error occured while deploying: %v", err)
		return
	}
	fmt.Println(res)
}
