package actions

import (
	"rft/internal/com/client"
	"rft/internal/op"
	"rft/internal/ui"
)

// Copy create a new the copy operation
func Copy(args map[string]interface{}) error {
	opArgs := op.OperationArgs(args)
	operation, err := client.CreateOperation(op.OpCopy, opArgs)

	if err != nil {
		return err
	}

	return ui.AddOperation(operation)
}
