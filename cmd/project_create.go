package cmd

import (
	"fmt"
	"github.com/spf13/viper"

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
	projectCreateCmd.LocalFlags().StringVar(&cfgFile, "config", ".cloudfauj.yml", "Project configuration file")
	projectCmd.AddCommand(projectCreateCmd)
}

func runProjectCreateCmd(cmd *cobra.Command, args []string) {
	/*
		1. Collect .cloudfauj.yml, server addr
		2. Initialize CF server client
		3. API call
		4. Display response
	*/
	serverAddr, _ := cmd.Flags().GetString("server-addr")
	fmt.Println(viper.Get("project"))

	fmt.Println(serverAddr)
	fmt.Println(args)
	fmt.Println("create called!")
}
