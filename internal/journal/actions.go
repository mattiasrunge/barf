package journal

import (
	"errors"

	"barf/internal/com/server"
	"barf/internal/op"
)

func create(opType op.OperationType, args op.OperationArgs) (*op.Operation, error) {
	operation := op.NewOperation(opType, args)

	if registerdStartHandler == nil {
		return operation, errors.New("No start handler registered")
	}

	e := entry{
		Operation: operation,
		Status:    op.NewStatus(),
	}

	err := addEntry(&e)

	if err != nil {
		return operation, err
	}

	err = server.OperationCreated(operation)

	if err != nil {
		UpdateOperationStatus(operation.ID, &op.OperationStatus{
			Finished: true,
			ExitCode: 255,
			Error:    err.Error(),
		})
		return operation, err
	}

	err = registerdStartHandler(operation)

	if err != nil {
		UpdateOperationStatus(operation.ID, &op.OperationStatus{
			Finished: true,
			ExitCode: 255,
			Error:    err.Error(),
		})
		return operation, err
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
			Error:    err.Error(),
		})
		return err
	}

	return nil
}

func list() ([]*op.Operation, error) {
	var operations []*op.Operation

	for _, e := range entries {
		operations = append(operations, e.Operation)
	}

	return operations, nil
}

func status(operationID op.OperationID) (*op.OperationStatus, error) {
	entry, err := getEntry(operationID)

	if err != nil {
		return nil, err
	}

	return entry.Status, errors.New("No such operation found")
}
