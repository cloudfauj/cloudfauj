package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var serverAddr string

var rootCmd = &cobra.Command{
	Use:   "cloudfauj",
	Short: "Deploy Apps to your cloud without managing infrastructure",
	Long: `
    CloudFauj enables developers to deploy applications in their own Cloud
    without having to manually provision or manage the infrastructure for them.

    If you've just installed cloudfauj, you can start by launching the server.

    If your cloudfauj server is already up & running, you can get started with
    deploying your applications.`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	envCmd.AddCommand(envCreateCmd, envDestroyCmd, envListCmd)
	deploymentCmd.AddCommand(deploymentStatusCmd, deploymentLogsCmd, deploymentListCmd)

	rootCmd.PersistentFlags().StringVar(
		&serverAddr,
		"server-addr",
		"http://127.0.0.1:6200",
		"HTTP address of Cloudfauj Server, including the Scheme",
	)
	rootCmd.AddCommand(serverCmd, envCmd, deployCmd, deploymentCmd)

	// prevent error message showing up twice
	rootCmd.SilenceErrors = true
	// Prevent usage from showing up when the command logic returns an error
	rootCmd.SilenceUsage = true
}

// initConfig loads configuration into viper from the given file.
// Because file can differ based on the command invoked, this func
// is invoked by individual command runners.
func initConfig(file string) {
	viper.SetConfigFile(file)
	if err := viper.ReadInConfig(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error while reading configuration: %v", err)
		os.Exit(1)
	}
}
