package actions

import (
	"barf/internal/com/client"
	"barf/internal/op"
	"barf/internal/ui"
)

// Dummy creates a new dummy operation
func Dummy(args map[string]interface{}) error {
	opArgs := op.OperationArgs{
		"iterations": args["iterations"],
	}

	operation, err := client.CreateOperation(op.OpDummy, opArgs)

	if err != nil {
		return err
	}

	return ui.AddOperation(operation)
}
