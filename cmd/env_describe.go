package cmd

import (
	"fmt"
	"github.com/cloudfauj/cloudfauj/api"
	"os"

	"github.com/spf13/cobra"
)

// envDescribeCmd represents the describe command
var envDescribeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describe an environment",
	Long:  "This command displays information about an environment.",
	Run:   runEnvDescribeCmd,
}

func runEnvDescribeCmd(cmd *cobra.Command, args []string) {
	serverAddr, _ := cmd.Flags().GetString("server-addr")
	apiClient := api.NewClient(serverAddr)

	env, err := apiClient.GetEnvironment(args[0])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "an error occured while fetching environment info: %v", err)
		return
	}
	fmt.Println(env)
}
