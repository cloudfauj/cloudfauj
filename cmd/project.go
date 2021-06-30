package cmd

import (
	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage Projects",
	Long: `Allows managing projects.

A Project represents 1 VCS repository where you write your business logic.
Each project has its own .cloudfauj.yml file that describes the applications
it contains and their infrastructure requirements. This file must be committed
in the root directory of the repo.

To get started, you can create a new project.`,
}
