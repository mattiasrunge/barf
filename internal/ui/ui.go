package ui

import (
	"barf/internal/com/client"
	"barf/internal/op"
	"errors"
	"os"
	"sync"

	"github.com/mattiasrunge/goterminal"
)

type operationWithStatus struct {
	operation *op.Operation
	status    *op.OperationStatus
}

var wg sync.WaitGroup
var operations []*operationWithStatus

func onOperationStatus(opID op.OperationID, status *op.OperationStatus) {
	for _, o := range operations {
		if o.operation.ID == opID && !o.status.Finished {
			o.status = status

			update()

			if status.Finished {
				wg.Done()
			}
		}
	}
}

// AddOperation adds operation to monitor
func AddOperation(operation *op.Operation) error {
	if writer == nil {
		return errors.New("Start has not been called")
	}

	operations = append(operations, &operationWithStatus{
		operation: operation,
		status:    &op.OperationStatus{},
	})

	wg.Add(1)

	update()

	return nil
}

// Start monitoring operations
func Start(w int) {
	writer = goterminal.New(os.Stdout)

	if w == 0 {
		width, _ = writer.GetTermDimensions()
	} else {
		width = w
	}

	client.OnOperationStatus(onOperationStatus)
}

// Wait for all operations to complete
func Wait() int {
	wg.Wait()
	writer.Reset()

	return 0
}
