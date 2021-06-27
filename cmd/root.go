package cmd

import (
	"github.com/spf13/cobra"
)

//var cfgFile string

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

//func init() {
//cobra.OnInitialize(initConfig)

// Here you will define your flags and configuration settings.
// Cobra supports persistent flags, which, if defined here,
// will be global for your application.

//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cloudfauj.yaml)")

// Cobra also supports local flags, which will only run
// when this action is called directly.
//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
//}

// initConfig reads in config file and ENV variables if set.
//func initConfig() {
//	if cfgFile != "" {
//		// Use config file from the flag.
//		viper.SetConfigFile(cfgFile)
//	} else {
//		// Find home directory.
//		home, err := homedir.Dir()
//		cobra.CheckErr(err)
//
//		// Search config in home directory with name ".cloudfauj" (without extension).
//		viper.AddConfigPath(home)
//		viper.SetConfigName(".cloudfauj")
//	}
//
//	viper.AutomaticEnv() // read in environment variables that match
//
//	// If a config file is found, read it in.
//	if err := viper.ReadInConfig(); err == nil {
//		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
//	}
//}
