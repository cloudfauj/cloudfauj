package cmd

import (
	"errors"
	"fmt"
	"github.com/cloudfauj/cloudfauj/api"
	"github.com/cloudfauj/cloudfauj/server"
	"github.com/spf13/cobra"
)

var tfApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Run terraform apply over infrastructure",
	Long: `
    This command runs terraform apply over a specified Cloudfauj component.

    It is most useful when changes are made to cloudfauj-managed TF configuration
    and need to be applied to achieve the desired state of infrastructure.

        cloudfauj tf apply --domain example.com

    Similarly, changes can be applied to an environment:

        cloudfauj tf apply --env staging

    NOTE: This feature currently has a limitation.
    It only applies the component specified and not its dependent infrastructure
    which may reside as separate TF projects.

    For eg- Running apply over a domain changes its infra, but not the
    environment(s) relying on it. If some change in the domain config affects
    its dependent envs, a separate apply needs to be run over the envs.

    Note that apply is currently not supported for applications.`,
	RunE: runTfApplyCmd,
}

func init() {
	f := tfApplyCmd.Flags()
	f.String("domain", "", "A domain registered with Cloudfauj")
	f.String("env", "", "An environment managed by Cloudfauj")
}

func runTfApplyCmd(cmd *cobra.Command, args []string) error {
	var eventsCh <-chan *server.Event

	apiClient, err := api.NewClient(serverAddr)
	if err != nil {
		return err
	}

	f := cmd.Flags()
	domain, _ := f.GetString("domain")
	env, _ := f.GetString("env")

	if domain != "" {
		eventsCh, err = apiClient.TFApplyDomain(domain)
	} else if env != "" {
		eventsCh, err = apiClient.TFApplyEnv(env)
	} else {
		return errors.New("either domain or environment must be passed to this command")
	}

	if err != nil {
		return err
	}
	for e := range eventsCh {
		if e.Err != nil {
			return e.Err
		}
		fmt.Println(e.Msg)
	}
	return nil
}
