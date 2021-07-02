package cmd

import (
	"fmt"
	"github.com/cloudfauj/cloudfauj/api"
	"github.com/spf13/cobra"
)

var envDestroyCmd = &cobra.Command{
	Use:   "destroy [flags] ENV",
	Short: "Destroy an Environment",
	Long: `
    This command lets you destroy an environment managed by Cloudfauj.

    It kills all running applications, cancels deployments and destroys all infrastructure
    of the environment. After destruction, the environment doesn't cost you money anymore.

    This command is idempotent and does nothing if the specified environment doesn't exist.`,
	Args:    cobra.ExactArgs(1),
	RunE:    runEnvDestroyCmd,
	Example: "cloudfauj env destroy staging",
}

func runEnvDestroyCmd(cmd *cobra.Command, args []string) error {
	apiClient, err := api.NewClient(serverAddr)
	if err != nil {
		return err
	}

	fmt.Printf("Destroying environment %s\n", args[0])
	if err := apiClient.DestroyEnvironment(args[0]); err != nil {
		return err
	}

	fmt.Println("Done")
	return nil
}
