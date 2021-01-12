package runner

import (
	"errors"

	"rft/internal/op"
)

func abort(operationID op.OperationID) error {
	runner := getRunner(operationID)

	if runner == nil {
		return errors.New("No such runner found")
	}

	return runner.Abort()
}
