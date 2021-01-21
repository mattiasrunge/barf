package actions

import (
	"barf/internal/com/client"
	"barf/internal/op"
	"barf/internal/ui"
)

// Pull create a new the pull operation
func Pull(args map[string]interface{}) error {
	opArgs := op.OperationArgs(args)
	operation, err := client.CreateOperation(op.OpPull, opArgs)

	if err != nil {
		return err
	}

	return ui.AddOperation(operation)
}
