package actions

import (
	"barf/internal/cli"
	"barf/internal/cli/ui"
	"barf/internal/com/client"
	"barf/internal/op"
)

// Dummy creates a new dummy operation
func Dummy(args map[string]interface{}) error {
	cli.Start()

	opArgs := op.OperationArgs{
		"iterations": args["iterations"],
	}

	operation, err := client.CreateOperation(op.OpDummy, opArgs)

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
