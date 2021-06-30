package cmd

import (
	"fmt"
	"github.com/cloudfauj/cloudfauj/api"
	"github.com/spf13/cobra"
	"os"
)

// projectListCmd represents the list command
var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	Long:  "This command lists all Projects managed by Cloudfauj.",
	Run:   runProjectListCmd,
}

func runProjectListCmd(cmd *cobra.Command, args []string) {
	serverAddr, _ := cmd.Flags().GetString("server-addr")
	apiClient := api.NewClient(serverAddr)

	projects, err := apiClient.ListProjects()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "an error occured while list all projects: %v", err)
		return
	}
	for p := range projects {
		fmt.Println(p)
	}
}
