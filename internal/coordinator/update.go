package coordinator

import (
	"errors"

	"barf/internal/com/server"
	"barf/internal/op"
)

// UpdateOperationStatus sets the operation status
func UpdateOperationStatus(operationID op.OperationID, status *op.OperationStatus) error {
	e := getEntry(operationID)

	if e == nil {
		return errors.New("Could not find any entry for operation with id " + string(operationID) + " to report status on!")
	}

	err := e.UpdateStatus(status)

	if err != nil {
		return err
	}

	if e.Status.Finished {
		removeEntry(e)
	}

	return server.OperationStatus(operationID, e.Status)
}
