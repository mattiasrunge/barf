package actions

import (
	"barf/internal/cli"
	"barf/internal/cli/ui"
	"barf/internal/com/client"
	"barf/internal/op"
)

// Pull create a new the pull operation
func Pull(args map[string]interface{}) error {
	cli.Start()

	opArgs := op.OperationArgs(args)
	operation, err := client.CreateOperation(op.OpPull, opArgs)

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
