package actions

import (
	"barf/internal/cli"
	"barf/internal/cli/ui"
	"barf/internal/com/client"
	"barf/internal/op"
)

// Move create a new the move operation
func Move(args map[string]interface{}) error {
	cli.Start()

	opArgs := op.OperationArgs(args)
	operation, err := client.CreateOperation(op.OpMove, opArgs)

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
