package coordinator

import (
	"errors"

	"barf/internal/com/server"
	"barf/internal/journal"
	"barf/internal/op"
)

func create(opType op.OperationType, args op.OperationArgs) (*op.Operation, error) {
	index, err := getNextIndex()

	if err != nil {
		return nil, err
	}

	operation := op.NewOperation(opType, args, index)

	e, err := journal.NewJournalEntry(operation)

	addEntry(e)

	_ = server.OperationCreated(operation)

	err = registerdStartHandler(operation)

	if err != nil {
		UpdateOperationStatus(operation.ID, &op.OperationStatus{
			Finished: true,
			ExitCode: 255,
			Message:  err.Error(),
		})

		return nil, err
	}

	return operation, nil
}

func abort(operationID op.OperationID) error {
	if registerdAbortHandler == nil {
		return errors.New("No abort handler registered")
	}

	err := registerdAbortHandler(operationID)

	if err != nil {
		UpdateOperationStatus(operationID, &op.OperationStatus{
			Finished: true,
			ExitCode: 254,
			Message:  err.Error(),
		})

		return err
	}

	return nil
}

func list() ([]*op.Operation, error) {
	var operations []*op.Operation

	for _, e := range getEntries() {
		operations = append(operations, e.Operation)
	}

	return operations, nil
}

func status(operationID op.OperationID) (*op.OperationStatus, error) {
	entry := getEntry(operationID)

	if entry == nil {
		return nil, errors.New("No entry for operation with id " + string(operationID) + " found!")
	}

	return entry.Status, nil
}
