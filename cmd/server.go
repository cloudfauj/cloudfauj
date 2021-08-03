package cmd

import (
	"database/sql"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/cloudfauj/cloudfauj/infrastructure"
	"github.com/cloudfauj/cloudfauj/server"
	"github.com/cloudfauj/cloudfauj/state"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"path"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Launch a Cloudfauj Server",
	Long: `
    This command starts a Cloudfauj Server that provisions infrastructure and
    manages environments & applications.

    The server sets up a local data directory if it doesn't already exist to manage
    all state.
    It is extremely important that you keep this dir continuously backed up to avoid
    losing track of all infrastructure.`,
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

	log := logrus.New()
	// todo: further configuration of logger

	if err := setupDataDir(log, srvCfg.DataDir); err != nil {
		return fmt.Errorf("failed to setup server data directory: %v", err)
	}

	db, err := sql.Open("sqlite3", srvCfg.DBFilePath())
	if err != nil {
		return fmt.Errorf("failed to open database connection: %v", err)
	}
	defer db.Close()

	storage := state.New(log, db)
	if err := storage.Migrate(cmd.Context()); err != nil {
		return fmt.Errorf("failed to run DB migrations: %v", err)
	}

	awsCfg, err := config.LoadDefaultConfig(cmd.Context())
	if err != nil {
		return fmt.Errorf("failed to setup AWS configuration: %v", err)
	}
	infra := infrastructure.New(
		log,
		ec2.NewFromConfig(awsCfg),
		ecs.NewFromConfig(awsCfg),
		iam.NewFromConfig(awsCfg),
		elasticloadbalancingv2.NewFromConfig(awsCfg),
		awsCfg.Region,
	)

	apiServer := server.New(&srvCfg, log, storage, infra)

	bindAddr := viper.GetString("bind_host") + ":" + viper.GetString("bind_port")
	log.WithFields(logrus.Fields{"bind_addr": bindAddr}).Info("Starting CloudFauj Server")

	if err := http.ListenAndServe(bindAddr, apiServer); err != nil {
		return fmt.Errorf("failed to start the server: %v", err)
	}
	return nil
}

func setupDataDir(log *logrus.Logger, dir string) error {
	// return if the data dir already exists
	_, err := os.Stat(dir)
	if err == nil {
		log.WithField("dir", dir).Debug("Data directory already exists")
		return nil
	}
	if !os.IsNotExist(err) {
		// unless the error is "dir not found", propagate the unexpected err
		return fmt.Errorf("failed to check if data directory already exists: %v", err)
	}

	log.WithField("dir", dir).Info("Setting up server data directory")

	var subDirs = []string{server.DBDir, server.DeploymentsDir}
	for _, sd := range subDirs {
		d := path.Join(dir, sd)
		if err := os.MkdirAll(d, 0755); err != nil {
			return fmt.Errorf("failed to create %s: %v", d, err)
		}
	}
	return nil
}
