package actions

import (
	"barf/internal/com/client"
	"barf/internal/op"
	"barf/internal/ui"
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
