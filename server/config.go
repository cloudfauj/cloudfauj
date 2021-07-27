package server

import "path"

type Config struct {
	// DataDir is the base directory inside which Cloudfauj server
	// stores all its data.
	// To restore Cloudfauj server on to a new server, restoring a
	// backup of this dir and running the server is enough.
	DataDir string `mapstructure:"data_dir"`
}

const (
	DeploymentsDir   = "deployments"
	LogFileBasename  = "logs.txt"
	DBDir            = "db"
	DBServerFilename = "server.db"
)

func (c *Config) DBFilePath() string {
	return path.Join(c.DataDir, DBDir, DBServerFilename)
}
