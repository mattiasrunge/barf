package journal

import (
	"barf/internal/com/server"
	"barf/internal/op"
)

// UpdateOperationStatus sets the operation status
func UpdateOperationStatus(operationID op.OperationID, status *op.OperationStatus) error {
	e, err := getEntry(operationID)

	if err != nil {
		return err
	}

	op.UpdateStatus(e.Status, status)

	err = writeEntry(e)

	if err != nil {
		return err
	}

	if e.Status.Finished {
		err = removeEntry(e)

		if err != nil {
			return err
		}
	}

	return server.OperationStatus(operationID, e.Status)
}
