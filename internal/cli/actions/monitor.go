package actions

import (
	"fmt"
	"sync"

	"barf/internal/cli"
	"barf/internal/cli/ui"
	"barf/internal/com/client"
	"barf/internal/op"
	"barf/internal/typeconv"
	"barf/internal/utils"
)

// Monitor starts monitoring all or specified operations
func Monitor(args map[string]interface{}) error {
	cli.Start()

	idxArray, _ := typeconv.ToArray(args["ids"])
	strIdx := typeconv.ToStringArray(idxArray)
	idx := typeconv.StringArray2IntArray(strIdx)

	operations, err := client.ListOperations()

	if err != nil {
		return err
	}

	for _, operation := range operations {
		if len(idx) == 0 || utils.IntArrayContains(idx, int(operation.Index)) {
			err = ui.AddOperation(operation)

			if err != nil {
				return err
			}
		}
	}

	if len(idx) == 0 {
		client.OnOperationCreated(func(operation *op.Operation) {
			err = ui.AddOperation(operation)

			if err != nil {
				fmt.Println(err)
			}
		})

		var wg sync.WaitGroup

		wg.Add(1)
		wg.Wait() // Wait until user aborts...
	}

	cli.Finish()

	return nil
}
