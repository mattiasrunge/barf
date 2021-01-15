package runners

import (
	"errors"

	"barf/internal/cmd"
	"barf/internal/config"
	"barf/internal/op"
)

type statushandler func(status *op.OperationStatus)

// Runner defines the executor interface
type Runner interface {
	init(operation *op.Operation)
	Start()
	Abort() error
	OperationID() op.OperationID
	OnStdout(cmd.LogHandler)
	OnStderr(cmd.LogHandler)
	OnStatus(statushandler)
}

// NewRunner creates an executor object of the correct type
func NewRunner(operation *op.Operation) (Runner, error) {
	var r Runner

	if operation.Type == op.OpDummy && !config.IsProduction() {
		r = &dummyRunner{}
	} else if operation.Type == op.OpCopy {
		r = &copyRunner{}
	} else if operation.Type == op.OpMove {
		r = &moveRunner{}
	} else {
		return nil, errors.New("Unknown operation type")
	}

	r.init(operation)

	return r, nil
}
