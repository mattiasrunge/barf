package actions

import (
	"barf/internal/cli"
	"barf/internal/cli/ui"
	"barf/internal/com/client"
	"barf/internal/op"
)

// Copy create a new the copy operation
func Copy(args map[string]interface{}) error {
	cli.Start()

	opArgs := op.OperationArgs(args)
	operation, err := client.CreateOperation(op.OpCopy, opArgs)

	if err != nil {
		return err
	}

	err = ui.AddOperation(operation)

	if err != nil {
		return err
	}

	cli.Finish()

	return nil
}
