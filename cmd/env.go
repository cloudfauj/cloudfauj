package cmd

import (
	"github.com/spf13/cobra"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage Environments",
	Long: `This command lets you work with Environments managed by Cloudfauj.

An environment is a group of applications that's logically isolated from other
environments. Some underlying infrastructure is created by CloudFauj in your
Cloud in order to support the env, so there's some cost attached to it.

The first thing you'd normally do is to create an environment.`,
}
