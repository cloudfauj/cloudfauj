package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/cloudfauj/cloudfauj/infrastructure"
	"github.com/cloudfauj/cloudfauj/server"
	"github.com/cloudfauj/cloudfauj/state"
	"github.com/hashicorp/terraform-exec/tfinstall"
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
	log := logger()

	// setup server configuration
	srvCfgFile, _ := cmd.Flags().GetString("config")
	initConfig(srvCfgFile)

	d := viper.GetString("data_dir")
	if d == "" {
		log.Warn("Server data directory not specified, using current directory")
		d, _ = os.Getwd()
	}
	srvCfg := server.NewConfig(d)

	// aws authentication
	log.Info("Validating AWS credentials")
	awsCfg, err := loadAWSConfig(cmd.Context())
	if err != nil {
		return err
	}

	// setup main data directory for server
	if err := setupDataDir(cmd.Context(), log, srvCfg); err != nil {
		return fmt.Errorf("failed to setup server data directory: %v", err)
	}

	// db setup
	db, err := sql.Open("sqlite3", srvCfg.DBFilePath())
	if err != nil {
		return fmt.Errorf("failed to open database connection: %v", err)
	}
	defer db.Close()

	storage := state.New(log, db)
	if err := storage.Migrate(cmd.Context()); err != nil {
		return fmt.Errorf("failed to run DB migrations: %v", err)
	}

	// TODO: Inject Terraform object from here
	//  Currently, we can't do this because we need to execute the shared tf object
	//  in different working dirs. This is only possible with the -chdir option,
	//  support for which is currently not merged into terraform-exec.
	//  See https://github.com/hashicorp/terraform-exec/pull/100.

	infra := &infrastructure.Infrastructure{
		Log:         log,
		Region:      awsCfg.Region,
		Ec2:         ec2.NewFromConfig(awsCfg),
		Ecs:         ecs.NewFromConfig(awsCfg),
		TFConfigDir: srvCfg.TerraformDir(),
		TFBinary:    path.Join(srvCfg.DataDir(), "terraform"),
	}
	apiServer := server.New(srvCfg, log, storage, infra)
	bindAddr := viper.GetString("bind_host") + ":" + viper.GetString("bind_port")

	log.WithFields(logrus.Fields{"bind_addr": bindAddr}).Info("Starting CloudFauj Server")
	if err := http.ListenAndServe(bindAddr, apiServer); err != nil {
		return fmt.Errorf("failed to start the server: %v", err)
	}
	return nil
}

func setupDataDir(ctx context.Context, log *logrus.Logger, srvCfg *server.Config) error {
	_, err := os.Stat(srvCfg.DataDir())
	if err == nil {
		log.WithField("dir", srvCfg.DataDir()).Info("Found data directory")
		return nil
	}
	if !os.IsNotExist(err) {
		// unless the error is "dir not found", propagate the unexpected err
		return fmt.Errorf("failed to check if data directory already exists: %v", err)
	}

	log.WithField("dir", srvCfg.DataDir()).Info("Setting up server data directory")
	subDirs := []string{
		srvCfg.DBDir(), srvCfg.DeploymentsDir(), srvCfg.TerraformDir(),
	}
	for _, sd := range subDirs {
		if err := os.MkdirAll(sd, 0755); err != nil {
			return fmt.Errorf("failed to create %s: %v", sd, err)
		}
	}
	return setupTerraform(ctx, log, srvCfg)
}

func setupTerraform(ctx context.Context, log *logrus.Logger, srvCfg *server.Config) error {
	log.WithField("version", srvCfg.TerraformVersion()).Info("Downloading Terraform")
	_, err := tfinstall.Find(
		ctx,
		tfinstall.ExactVersion(srvCfg.TerraformVersion(), srvCfg.DataDir()),
	)
	if err != nil {
		return fmt.Errorf("failed to locate Terraform binary: %s", err)
	}
	return nil
}

func logger() *logrus.Logger {
	l := logrus.New()
	l.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	return l
}

func loadAWSConfig(ctx context.Context) (aws.Config, error) {
	c, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return c, fmt.Errorf("failed to setup AWS configuration: %v", err)
	}
	if c.Region == "" {
		return c, fmt.Errorf("no AWS region specified")
	}
	if _, err := c.Credentials.Retrieve(ctx); err != nil {
		return c, fmt.Errorf("no AWS credentials supplied, cannot proceed")
	}
	return c, nil
}
