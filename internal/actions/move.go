package actions

import (
	"barf/internal/com/client"
	"barf/internal/op"
	"barf/internal/ui"
)

// Move create a new the move operation
func Move(args map[string]interface{}) error {
	opArgs := op.OperationArgs(args)
	operation, err := client.CreateOperation(op.OpMove, opArgs)

	if err != nil {
		return err
	}

	return ui.AddOperation(operation)
}
