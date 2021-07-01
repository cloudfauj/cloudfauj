package cmd

import (
	"fmt"
	"github.com/cloudfauj/cloudfauj/api"
	"os"

	"github.com/spf13/cobra"
)

var envDestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy an Environment",
	Long: `This command lets you destroy an environment managed by Cloudfauj.

It kills all running applications, cancels deployments and destroys all infrastructure
of the environment. After destruction, the environment doesn't cost you money anymore.

This command is idempotent and does nothing if the specified environment doesn't exist.

Below example destroys an environment named "staging":
    cloudfauj env destroy staging`,
	Run: runEnvDestroyCmd,
}

func runEnvDestroyCmd(cmd *cobra.Command, args []string) {
	serverAddr, _ := cmd.Flags().GetString("server-addr")
	apiClient := api.NewClient(serverAddr)

	fmt.Printf("Destroying environment %s\n", args[0])
	if err := apiClient.DestroyEnvironment(args[0]); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "an error occured while destroying the environment: %v", err)
		return
	}
	fmt.Println("Done")
}
