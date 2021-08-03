package cmd

import (
	"fmt"
	"github.com/cloudfauj/cloudfauj/api"
	"github.com/spf13/cobra"
)

var appDestroyCmd = &cobra.Command{
	Use:   "destroy --env ENV [flags] APP",
	Args:  cobra.ExactArgs(1),
	Short: "Destroy an application",
	Long: `
    This command lets you stop an application and destroy all infrastructure
    that was provisioned for it in a specific environment.`,
	RunE:    runAppDestroyCmd,
	Example: "cloudfauj app destroy --env staging demo-server",
}

func init() {
	appDestroyCmd.Flags().String("env", "", "The environment to destroy the app from")
	_ = appDestroyCmd.MarkFlagRequired("env")
}

func runAppDestroyCmd(cmd *cobra.Command, args []string) error {
	env, _ := cmd.Flags().GetString("env")
	apiClient, err := api.NewClient(serverAddr)
	if err != nil {
		return err
	}
	if err := apiClient.DestroyApp(args[0], env); err != nil {
		return err
	}
	fmt.Printf("Destroyed %s from %s\n", args[0], env)
	return nil
}
