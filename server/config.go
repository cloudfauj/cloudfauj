package server

import (
	"fmt"
	"path"
)

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

const (
	TerraformDir                = "infrastructure"
	TerraformVersion            = "1.0.4"
	TerraformConfFile           = "terraform.tf"
	TerraformAwsProviderVersion = "3.54.0"
)

const tfConfig = `terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "%s"
    }
  }
}

provider "aws" {
  region = "%s"
}
`

func TfConfig(region string) string {
	return fmt.Sprintf(tfConfig, TerraformAwsProviderVersion, region)
}

func (c *Config) DBFilePath() string {
	return path.Join(c.DataDir, DBDir, DBServerFilename)
}
