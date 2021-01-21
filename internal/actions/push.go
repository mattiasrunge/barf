package actions

import (
	"barf/internal/com/client"
	"barf/internal/op"
	"barf/internal/ui"
)

// Push create a new the push operation
func Push(args map[string]interface{}) error {
	opArgs := op.OperationArgs(args)
	operation, err := client.CreateOperation(op.OpPush, opArgs)

	if err != nil {
		return err
	}

	return ui.AddOperation(operation)
}
