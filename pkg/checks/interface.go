package checks

import (
	"context"
	"io"
)

type Check interface {
	Name() string
	Run(ctx context.Context, commandLog io.Writer) Result
}

type Result struct {
	Errors    []error
	PostCheck []string
}
