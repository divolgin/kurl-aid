package checks

import (
	"context"
	"io"
	"os/exec"

	"github.com/divolgin/kurl-aid/pkg/log"
	"github.com/pkg/errors"
)

type Kubectl struct {
}

func (c Kubectl) Name() string {
	return "kubectl"
}

// TODO: Detect cluster type. This works for kurl
func (c Kubectl) Run(ctx context.Context, commandLog io.Writer) (result Result) {
	result.Errors = []error{}
	result.PostCheck = []string{}

	log.Printf("Checking if kubectl exists")
	_, err := exec.LookPath("kubectl")
	if err != nil {
		result.Errors = append(result.Errors, errors.Wrap(err, "lookup path"))
	}

	return
}
