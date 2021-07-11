package cmd

import (
	"fmt"
	"github.com/cloudfauj/cloudfauj/server"
	"github.com/cloudfauj/cloudfauj/state"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Launch a Cloudfauj Server",
	Long: `
    This command starts a Cloudfauj Server that carries out tasks such
    as Deployments when requested.
    The server takes care of provisioning and managing all infrastructure required
    by applications.`,
	RunE: runServerCmd,
}

func init() {
	serverCmd.Flags().String("config", "", "Server configuration file")
	_ = serverCmd.MarkFlagRequired("config")
}

func runServerCmd(cmd *cobra.Command, args []string) error {
	configFile, _ := cmd.Flags().GetString("config")
	initConfig(configFile)

	var srvCfg server.Config
	if err := viper.Unmarshal(&srvCfg); err != nil {
		return fmt.Errorf("failed to parse server configuration: %v", err)
	}

	// todo: setup local data storage directory

	log := logrus.New()
	// todo: further configuration of logger

	storage := state.New(log)
	apiServer := server.New(&srvCfg, log, storage)

	bindAddr := viper.GetString("bind_host") + ":" + viper.GetString("bind_port")
	log.WithFields(logrus.Fields{"bind_addr": bindAddr}).Info("Starting CloudFauj Server")

	if err := http.ListenAndServe(bindAddr, apiServer); err != nil {
		return fmt.Errorf("failed to start the server: %v", err)
	}

	return nil
}
