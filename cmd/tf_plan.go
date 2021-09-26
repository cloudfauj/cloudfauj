package cmd

import (
	"errors"
	"fmt"
	"github.com/cloudfauj/cloudfauj/api"
	"github.com/cloudfauj/cloudfauj/server"
	"github.com/spf13/cobra"
)

var tfPlanCmd = &cobra.Command{
	Use:   "plan",
	Short: "Run Terraform plan over infrastructure",
	Long: `
    This command runs Terraform plan over a specified Cloudfauj component.

    It is most useful when changes are made to cloudfauj-managed TF configuration
    or cloud resources manually and a diff is needed between the desired and actual
    state of the infrastructure.

    For eg- After making changes to a domain's TF configuration, the plan command
    can be invoked to produce possible changes to the domain infra.

        cloudfauj tf plan --domain example.com

    Similarly, plans can be generated for an environment or a particular app
    in it.

        cloudfauj tf plan --env staging
        cloudfauj tf plan --env staging --app demo_api

    NOTE: This feature currently has a limitation.
    It only plans the component specified and not its dependent infrastructure
    which may reside as a separate TF project.

    For eg- Running plan over a domain shows the diff for it, but not for the
    environment(s) relying on it. If some change in the domain config affects
    its dependent envs, a separate plan needs to be run over the envs.`,
	RunE: runTfPlanCmd,
}

func init() {
	f := tfPlanCmd.Flags()
	f.String("domain", "", "A domain registered with Cloudfauj")
	f.String("env", "", "An environment managed by Cloudfauj")
	f.String("app", "", "An application managed by Cloudfauj (env must be specified)")
}

func runTfPlanCmd(cmd *cobra.Command, args []string) error {
	var eventsCh <-chan *server.Event

	apiClient, err := api.NewClient(serverAddr)
	if err != nil {
		return err
	}

	f := cmd.Flags()
	domain, _ := f.GetString("domain")
	env, _ := f.GetString("env")
	app, _ := f.GetString("app")

	if domain != "" {
		eventsCh, err = apiClient.TFPlanDomain(domain)
	} else if env != "" {
		if app != "" {
			eventsCh, err = apiClient.TFPlanApp(env, app)
		} else {
			eventsCh, err = apiClient.TFPlanEnv(env)
		}
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
