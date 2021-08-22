package server

import "path"

// Configuration passed to a Cloudfauj server.
// It dictates where and how the server state is organized.
type Config struct {
	// Base directory inside which Cloudfauj server stores all its data.
	// To restore Cloudfauj server on to a new server, restoring a
	// backup of this dir and running the server is enough.
	dataDir string

	// The directory inside base containing all data about deployments.
	deploymentsDir string

	// Name given to every deployment log file
	logfileName string

	// The directory inside base containing the database file(s).
	dbDir string

	// Database file name
	dbFilename string

	// Directory inside base containing all terraform configurations
	terraformDir string

	// Name of the main Terraform config file.
	// The value of this is always "terraform.tf".
	terraformConfigFile string

	// The terraform version the server works with
	terraformVersion string
}

// NewConfig returns a new Server Configuration
func NewConfig(dataDir string) *Config {
	return &Config{
		dataDir:             dataDir,
		deploymentsDir:      "deployments",
		logfileName:         "logs.txt",
		dbDir:               "db",
		dbFilename:          "server.db",
		terraformDir:        "infrastructure",
		terraformConfigFile: "terraform.tf",
		terraformVersion:    "1.0.5",
	}
}

func (c *Config) DataDir() string {
	return c.dataDir
}

// DBDir returns the exact path of the directory containing database file(s)
func (c *Config) DBDir() string {
	return path.Join(c.DataDir(), c.dbDir)
}

// DBFilePath returns the exact path of main database file
func (c *Config) DBFilePath() string {
	return path.Join(c.DBDir(), c.dbFilename)
}

// DeploymentsDir returns the exact path of directory containing
// data of all deployments.
func (c *Config) DeploymentsDir() string {
	return path.Join(c.DataDir(), c.deploymentsDir)
}

// TerraformDir returns the exact path of directory containing
// all terraform infrastructure configurations.
func (c *Config) TerraformDir() string {
	return path.Join(c.DataDir(), c.terraformDir)
}

// TerraformVersion returns the version of terraform the server
// works with.
func (c *Config) TerraformVersion() string {
	return c.terraformVersion
}
