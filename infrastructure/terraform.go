package infrastructure

import (
	"fmt"
	"github.com/hashicorp/terraform-exec/tfexec"
	"strings"
	"text/template"
)

const terraformAwsProviderVersion = "3.55.0"

func (i *Infrastructure) NewTerraform(workDir string) (*tfexec.Terraform, error) {
	tf, err := tfexec.NewTerraform(workDir, i.TFBinary)
	if err != nil {
		return nil, fmt.Errorf("failed to create new terraform object: %s", err)
	}

	// Pass the server process' environment variables to Terraform process
	tf.SetEnv(nil)
	// Set logging
	tf.SetLogger(i.Log)
	tf.SetStderr(i.Log.Out)
	tf.SetStdout(i.Log.Out)

	return tf, nil
}

func (i *Infrastructure) tfCoreConfig() string {
	var b strings.Builder
	t := template.Must(template.New("").Parse(tfCoreConfigTpl))
	data := map[string]interface{}{
		"aws_region":           i.Region,
		"aws_provider_version": terraformAwsProviderVersion,
	}
	t.Execute(&b, data)
	return b.String()
}
