package checks

import (
	"context"
	"io"
	"os"

	"github.com/divolgin/kurl-aid/pkg/log"
	"github.com/divolgin/kurl-aid/pkg/system"
	"github.com/divolgin/kurl-aid/pkg/ui"
	"github.com/pkg/errors"
)

type KubeConfig struct {
}

func (c KubeConfig) Name() string {
	return "kubeconfig"
}

func (c KubeConfig) Run(ctx context.Context, commandLog io.Writer) (result Result) {
	result.Errors = []error{}
	result.PostCheck = []string{}

	envPath := os.Getenv("KUBECONFIG")
	if envPath != "/etc/kubernetes/admin.conf" {
		prompt := `
	$KUBECONFIG should be set to /etc/kubernetes/admin.conf.

	Would you like update $KUBECONFIG variable in /etc/profile?`
		choice, err := ui.Prompt(prompt)
		if err != nil {
			result.Errors = append(result.Errors, errors.Wrap(err, "user prompt"))
			return
		}

		if choice == ui.No {
			result.Errors = append(result.Errors, errors.New("$KUBECONFIG is not set in /etc/profile"))
		} else {
			log.Printf("Adding KUBECONFIG to /etc/profile")
			output, err := system.CombinedOutput(commandLog, "bash", "-c", `echo "export KUBECONFIG=/etc/kubernetes/admin.conf" >> /etc/profile`)
			if err != nil {
				result.Errors = append(result.Errors, errors.Wrapf(err, "add KUBECONFIG: %s", output))
			} else {
				result.PostCheck = append(result.PostCheck, `bash -l`)
			}
		}
	}

	log.Printf("Checking permissions on /etc/kubernetes/admin.conf")
	file, err := os.Open("/etc/kubernetes/admin.conf")
	if err == nil {
		file.Close()
		return
	}

	if !os.IsPermission(err) {
		result.Errors = append(result.Errors, errors.Wrap(err, "open /etc/kubernetes/admin.conf"))
		return
	}

	prompt := `
	This user does not have read access to /etc/kubernetes/admin.conf. This is required to run kubectl commands.

	Would you like give this user read access to /etc/kubernetes/admin.conf?`
	choice, err := ui.Prompt(prompt)
	if err != nil {
		result.Errors = append(result.Errors, errors.Wrap(err, "user prompt"))
		return
	}

	if choice == ui.No {
		result.Errors = append(result.Errors, errors.New("no access to /etc/kubernetes/admin.conf"))
		return
	}

	log.Printf("Setting permissions on /etc/kubernetes/admin.conf")
	output, err := system.CombinedOutput(commandLog, "sudo", "chmod", "a+r", "/etc/kubernetes/admin.conf")
	if err != nil {
		result.Errors = append(result.Errors, errors.Wrapf(err, "update permissions: %s", output))
	}

	return
}
