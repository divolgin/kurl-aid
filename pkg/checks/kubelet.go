package checks

import (
	"context"
	"io"

	"github.com/divolgin/kurl-aid/pkg/log"
	"github.com/divolgin/kurl-aid/pkg/system"
	"github.com/divolgin/kurl-aid/pkg/ui"
	"github.com/pkg/errors"
)

type Kubelet struct {
}

func (c Kubelet) Name() string {
	return "kubelet"
}

func (c Kubelet) Run(ctx context.Context, commandLog io.Writer) (result Result) {
	result.Errors = []error{}
	result.PostCheck = []string{}

	log.Printf("Checking if kubelet is running")
	err := c.checkIsRunning(ctx, commandLog)
	if err != nil {
		result.Errors = append(result.Errors, errors.Wrap(err, "kubelet running"))
		return
	}

	return
}

func (c Kubelet) checkIsRunning(ctx context.Context, commandLog io.Writer) error {
	status, err := system.GetServiceStatus("kubelet", commandLog)
	if err != nil {
		return errors.Wrap(err, "get service status")
	}

	if status == system.ServiceNoService {
		return errors.New("kubelet service is missing")
	}

	if status == system.ServiceRunning {
		return nil
	}

	prompt := `
	Kubelet is not running.

	Would you like to start it?`
	choice, err := ui.Prompt(prompt)
	if err != nil {
		return errors.Wrap(err, "user prompt")
	}

	if choice == ui.No {
		return errors.New("kubelet is not running")
	}

	log.Printf("Starting kubelet")
	err = system.StartService("kubelet", commandLog)
	if err != nil {
		return errors.Wrap(err, "start kubelet")
	}

	return nil
}
