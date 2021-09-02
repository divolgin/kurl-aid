package system

import (
	"io"
	"os/exec"

	"github.com/pkg/errors"
)

type ServiceStatus string

const (
	ServiceRunning    ServiceStatus = "running"
	ServiceDead                     = "dead"
	ServiceNotRunning               = "not running"
	ServiceNoService                = "no service"
	ServiceUnknown                  = "unknown"
)

func GetServiceStatus(service string, commandLog io.Writer) (ServiceStatus, error) {
	_, err := CombinedOutput(commandLog, "systemctl", "status", service)
	if err == nil {
		return ServiceRunning, nil
	}

	if exiterr, ok := err.(*exec.ExitError); ok {
		switch exiterr.ExitCode() {
		case 1, 2:
			return ServiceDead, nil
		case 3:
			return ServiceNotRunning, nil
		case 4:
			return ServiceNoService, nil
		}
	}

	return ServiceUnknown, errors.Wrap(err, "run systemctl")
}

func StartService(service string, commandLog io.Writer) error {
	output, err := CombinedOutput(commandLog, "sudo", "systemctl", "start", service)
	if err != nil {
		return errors.Wrapf(err, "start service: %s", output)
	}
	return nil
}

func StopService(service string, commandLog io.Writer) error {
	output, err := CombinedOutput(commandLog, "sudo", "systemctl", "stop", service)
	if err != nil {
		return errors.Wrapf(err, "stop service: %s", output)
	}
	return nil
}

func DisableService(service string, commandLog io.Writer) error {
	output, err := CombinedOutput(commandLog, "sudo", "systemctl", "disable", service)
	if err != nil {
		return errors.Wrapf(err, "disable service: %s", output)
	}
	return nil
}
