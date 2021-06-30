package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "cloudfauj",
	Short: "Deploy Apps to your cloud without managing infrastructure",
	Long: `CloudFauj enables you to deploy your applications in your own Cloud
without having to manually provision the infrastructure to support it.

Launch the Server so it can start accepting and executing deployment jobs.
Use other commands such as deploy to interact with the server and carry out tasks.
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	projectCmd.AddCommand(projectCreateCmd, projectListCmd)
	envCmd.AddCommand(envCreateCmd, envDescribeCmd, envListCmd)
	deploymentCmd.AddCommand(deploymentStatusCmd, deploymentLogsCmd)

	rootCmd.PersistentFlags().String("server-addr", "http://127.0.0.1:6200", "Cloudfauj Server address")
	rootCmd.AddCommand(serverCmd, projectCmd, envCmd, deployCmd, deploymentCmd)
}

// initConfig loads configuration into viper from the given file.
// Because file can differ based on the command invoked, this func
// is invoked by individual command runners.
func initConfig(file string) {
	viper.SetConfigFile(file)
	if err := viper.ReadInConfig(); err != nil {
		_ = fmt.Errorf("error while reading configuration: %v", err)
	}
}
