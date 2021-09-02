package checks

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"os/exec"
	"strings"

	"github.com/divolgin/kurl-aid/pkg/log"
	"github.com/divolgin/kurl-aid/pkg/system"
	"github.com/divolgin/kurl-aid/pkg/ui"
	"github.com/pkg/errors"
)

type SELinux struct {
}

type SEStatus struct {
	Status      string
	CurrentMode string
	ConfigMode  string
}

func (c SELinux) Name() string {
	return "selinux"
}

func (c SELinux) Run(ctx context.Context, commandLog io.Writer) (result Result) {
	result.Errors = []error{}
	result.PostCheck = []string{}

	_, err := exec.LookPath("sestatus")
	if err != nil {
		if err, ok := err.(*exec.Error); ok && err.Unwrap() == exec.ErrNotFound {
			return
		}
		result.Errors = append(result.Errors, errors.Wrap(err, "lookup path"))
		return
	}

	sestatus, err := c.getSEStatus(commandLog)
	if err != nil {
		result.Errors = append(result.Errors, errors.Wrap(err, "get sestatus"))
		return
	}

	if sestatus.CurrentMode == "enforcing" {
		prompt := `
		SELinux is currently in enforcing mode. This is incompatible with kubernetes.
	
		Would you like to change it to permissive?`
		choice, err := ui.Prompt(prompt)
		if err != nil {
			result.Errors = append(result.Errors, errors.Wrap(err, "user prompt"))
			return
		}

		if choice == ui.No {
			result.Errors = append(result.Errors, errors.New("selinux is enforcing"))
			return
		}

		log.Printf("Setting SELinux current mode to permissive")
		err = c.setCurrentModePermissive(commandLog)
		if err != nil {
			result.Errors = append(result.Errors, errors.Wrap(err, "set current mode"))
		}
	}

	if sestatus.ConfigMode == "enforcing" {
		prompt := `
		SELinux's config is in enforcing mode. SELinux will become enforcing after system reboot. This is incompatible with kubernetes.
	
		Would you like to permanently change it to permissive?`
		choice, err := ui.Prompt(prompt)
		if err != nil {
			result.Errors = append(result.Errors, errors.Wrap(err, "user prompt"))
			return
		}

		if choice == ui.No {
			result.Errors = append(result.Errors, errors.New("selinux config is enforcing"))
			return
		}

		log.Printf("Setting SELinux config mode to permissive")
		err = c.setConfigModePermissive(commandLog)
		if err != nil {
			result.Errors = append(result.Errors, errors.Wrap(err, "set config mode"))
		}
	}

	return
}

func (c SELinux) getSEStatus(commandLog io.Writer) (*SEStatus, error) {
	output, err := system.CombinedOutput(commandLog, "sestatus")
	if err != nil {
		return nil, errors.Wrapf(err, "run sestatus: %s", output)
	}

	result := &SEStatus{}

	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ":")
		if len(parts) != 2 {
			continue
		}
		switch parts[0] {
		case "SELinux status":
			result.Status = strings.TrimSpace(parts[1])
		case "Current mode":
			result.CurrentMode = strings.TrimSpace(parts[1])
		case "Mode from config file":
			result.ConfigMode = strings.TrimSpace(parts[1])
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, errors.Wrap(err, "read status")
	}

	return result, nil
}

func (c SELinux) setCurrentModePermissive(commandLog io.Writer) error {
	output, err := system.CombinedOutput(commandLog, "sudo", "setenforce", "0")
	if err != nil {
		return errors.Wrapf(err, "set current permissive: %s", output)
	}
	return nil
}

func (c SELinux) setConfigModePermissive(commandLog io.Writer) error {
	output, err := system.CombinedOutput(commandLog, "sudo", "sed", "-i", `s/^SELINUX=.*$/SELINUX=permissive/`, "/etc/selinux/config")
	if err != nil {
		return errors.Wrapf(err, "set config permissive: %s", output)
	}
	return nil
}
