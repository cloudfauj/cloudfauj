package cmd

import (
	"fmt"
	"github.com/cloudfauj/cloudfauj/api"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

var projectCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Project",
	Long: `This command lets you create a new project.
As a developer, this is the first thing you'd normally do after installing Cloudfauj.

Make sure a .cloudfauj.yml configuration is present in the current directory.
Below is an example command:

    # Registers the project & its applications with cloudfauj server
    cloudfauj project create`,
	Run: runProjectCreateCmd,
}

func init() {
	projectCreateCmd.Flags().String("config", ".cloudfauj.yml", "Project configuration file")
}

func runProjectCreateCmd(cmd *cobra.Command, args []string) {
	serverAddr, _ := cmd.Flags().GetString("server-addr")
	apiClient := api.NewClient(serverAddr)

	configFile, _ := cmd.Flags().GetString("config")
	initConfig(configFile)

	project := viper.GetString("project")
	apps := viper.GetStringMap("applications")

	fmt.Printf("Creating project %s\n", project)
	if err := apiClient.CreateProject(project, apps); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "an error occured while creating project: %v", err)
		return
	}
	fmt.Println("Done")
}
