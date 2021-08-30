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
    CloudFauj empowers developers to deploy applications in their own Cloud
    without having to manually provision or manage the infrastructure for them.

    If you've just installed cloudfauj, you can start by launching the server.

    If your cloudfauj server is already up & running, you can get started by
    creating an environment and deploying an application to it.`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	appCmd.AddCommand(appDestroyCmd)
	envCmd.AddCommand(envCreateCmd, envDestroyCmd, envListCmd)
	deploymentCmd.AddCommand(deploymentInfoCmd, deploymentLogsCmd, deploymentListCmd)
	domainCmd.AddCommand(domainAddCmd, domainDeleteCmd)

	rootCmd.PersistentFlags().StringVar(
		&serverAddr,
		"server-addr",
		"http://127.0.0.1:6200",
		"HTTP address of Cloudfauj Server, including the Scheme",
	)
	rootCmd.AddCommand(serverCmd, envCmd, appCmd, deployCmd, deploymentCmd, domainCmd, domainDeleteCmd)

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
