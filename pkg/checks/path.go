package checks

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/divolgin/kurl-aid/pkg/log"
	"github.com/divolgin/kurl-aid/pkg/system"
	"github.com/divolgin/kurl-aid/pkg/ui"
	"github.com/pkg/errors"
)

type Path struct {
}

func (c Path) Name() string {
	return "PATH"
}

func (c Path) Run(ctx context.Context, commandLog io.Writer) (result Result) {
	result.Errors = []error{}
	result.PostCheck = []string{}

	envPath := os.Getenv("PATH")
	pathDirs := strings.Split(envPath, ":")
	for _, pathDir := range pathDirs {
		if pathDir == "/usr/local/bin" {
			return
		}
	}

	prompt := `
	$PATH does not include the /usr/local/bin. This is required in order to use kubectl plugins.

	Would you like update $PATH variable in ~/.bash_profile?`
	choice, err := ui.Prompt(prompt)
	if err != nil {
		result.Errors = append(result.Errors, errors.Wrap(err, "user prompt"))
		return
	}

	if choice == ui.No {
		result.Errors = append(result.Errors, errors.New("$PATH does not include \"/usr/local/bin\""))
		return
	}

	log.Printf("Adding PATH to ~/.bash_profile")
	output, err := system.CombinedOutput(commandLog, "bash", "-c", `echo "PATH=\$PATH:/usr/local/bin" >> ~/.bash_profile`)
	if err != nil {
		result.Errors = append(result.Errors, errors.Wrapf(err, "add path: %s", output))
		return
	}

	output, err = system.CombinedOutput(commandLog, "bash", "-c", `echo "export PATH" >> ~/.bash_profile`)
	if err != nil {
		result.Errors = append(result.Errors, errors.Wrapf(err, "add export: %s", output))
		return
	}

	result.PostCheck = append(result.PostCheck, `bash -l`)

	return
}
