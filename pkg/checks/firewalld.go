package checks

import (
	"context"
	_ "embed"
	"io"

	"github.com/divolgin/kurl-aid/pkg/log"
	"github.com/divolgin/kurl-aid/pkg/system"
	"github.com/divolgin/kurl-aid/pkg/ui"
	"github.com/pkg/errors"
)

type FirewallD struct {
}

func (c FirewallD) Name() string {
	return "firewalld"
}

func (c FirewallD) Run(ctx context.Context, commandLog io.Writer) (result Result) {
	result.Errors = []error{}
	result.PostCheck = []string{}

	status, err := system.GetServiceStatus("firewalld", commandLog)
	if err != nil {
		result.Errors = append(result.Errors, errors.Wrap(err, "get service status"))
		return
	}

	if status != system.ServiceRunning {
		return
	}

	prompt := `
	Firewalld is running. This can prevent containers from communicating with each other and with external systems.

	Would you like to stop and disable it?`
	choice, err := ui.Prompt(prompt)
	if err != nil {
		result.Errors = append(result.Errors, errors.Wrap(err, "user prompt"))
		return
	}

	if choice == ui.No {
		result.Errors = append(result.Errors, errors.New("firewalld is running"))
		return
	}

	log.Printf("Stopping firewalld")
	err = system.StopService("firewalld", commandLog)
	if err != nil {
		result.Errors = append(result.Errors, errors.Wrap(err, "stop firewalld"))
	}

	log.Printf("Disabling firewalld")
	err = system.DisableService("firewalld", commandLog)
	if err != nil {
		result.Errors = append(result.Errors, errors.Wrap(err, "disable firewalld"))
	}

	return result
}
