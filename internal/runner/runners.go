package runner

import (
	"fmt"

	"barf/internal/journal"
	"barf/internal/op"
	"barf/internal/runner/runners"
)

var activeRunners []runners.Runner

func createRunner(operation *op.Operation) error {
	r, err := runners.NewRunner(operation)

	if err != nil {
		return err
	}

	r.OnStdout(func(line string) {
		fmt.Println("Stdout["+r.OperationID()+"]: ", line)
		// TODO: Write to file
		// TODO: Send to clients?
	})

	r.OnStderr(func(line string) {
		fmt.Println("Stderr["+r.OperationID()+"]: ", line)
		// TODO: Write to file
		// TODO: Send to clients?
	})

	r.OnStatus(func(status *op.OperationStatus) {
		err := journal.UpdateOperationStatus(r.OperationID(), status)

		if err != nil {
			fmt.Println(err)
		}

		if status.Finished {
			removeRunner(r.OperationID())
		}
	})

	activeRunners = append(activeRunners, r)

	go r.Start()

	return nil
}

func removeRunner(operationID op.OperationID) {
	for i, r := range activeRunners {
		if r.OperationID() == operationID {
			copy(activeRunners[i:], activeRunners[i+1:])
			activeRunners = activeRunners[:len(activeRunners)-1]

			return
		}
	}
}

func getRunner(operationID op.OperationID) runners.Runner {
	for _, r := range activeRunners {
		if r.OperationID() == operationID {
			return r
		}
	}

	return nil
}
