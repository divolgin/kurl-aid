package checks

import (
	"context"
	"io"
)

type IPTables struct {
}

func (c IPTables) Name() string {
	return "iptables"
}

func (c IPTables) Run(ctx context.Context, commandLog io.Writer) (result Result) {
	result.Errors = []error{}
	result.PostCheck = []string{}
	return result
}
