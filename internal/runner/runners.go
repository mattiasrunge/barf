package runner

import (
	"fmt"
	"sync"

	"barf/internal/coordinator"
	"barf/internal/op"
	"barf/internal/runner/runners"
)

var activeRunners []runners.Runner
var activeRunnersMu sync.Mutex

func createRunner(operation *op.Operation) error {
	r, err := runners.NewRunner(operation)

	if err != nil {
		return err
	}

	r.OnStdout(func(line string) {
		fmt.Println("Stdout["+r.OperationID()+"]: ", line)

		coordinator.WriteOperationStdout(r.OperationID(), line)
	})

	r.OnStderr(func(line string) {
		fmt.Println("Stderr["+r.OperationID()+"]: ", line)

		coordinator.WriteOperationStderr(r.OperationID(), line)
	})

	r.OnStatus(func(status *op.OperationStatus) {
		err := coordinator.UpdateOperationStatus(r.OperationID(), status)

		if err != nil {
			fmt.Println(err)
		}

		if status.Finished {
			removeRunner(r.OperationID())
		}
	})

	activeRunnersMu.Lock()
	activeRunners = append(activeRunners, r)
	defer activeRunnersMu.Unlock()

	go r.Start()

	return nil
}

func removeRunner(operationID op.OperationID) {
	activeRunnersMu.Lock()
	defer activeRunnersMu.Unlock()

	for i, r := range activeRunners {
		if r.OperationID() == operationID {
			copy(activeRunners[i:], activeRunners[i+1:])
			activeRunners = activeRunners[:len(activeRunners)-1]

			return
		}
	}
}

func getRunner(operationID op.OperationID) runners.Runner {
	activeRunnersMu.Lock()
	defer activeRunnersMu.Unlock()

	for _, r := range activeRunners {
		if r.OperationID() == operationID {
			return r
		}
	}

	return nil
}
