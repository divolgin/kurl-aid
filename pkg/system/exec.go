package system

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

func Exec(command string, args ...string) ([]byte, []byte, error) {
	cmd := exec.Command(command, args...)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()

	return stdout.Bytes(), stderr.Bytes(), errors.Wrap(err, "exec command")
}

func CombinedOutput(commandLog io.Writer, command string, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()

	_, _ = fmt.Fprintln(commandLog, command, strings.Join(args, " "))
	_, _ = fmt.Fprintln(commandLog, string(output))

	return output, err
}
