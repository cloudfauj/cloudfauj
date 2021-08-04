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

    Before destroying an environment, you must ensure that no applications
    are deployed to it. This will be automated in future.

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
	fmt.Printf("Destroying %s\n\n", args[0])
	eventsCh, err := apiClient.DestroyEnvironment(args[0])
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
