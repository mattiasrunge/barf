package coordinator

import (
	"errors"

	"barf/internal/com/server"
	"barf/internal/journal"
	"barf/internal/op"
)

// Initialize reads active entries from disk
func Initialize() error {
	if registerdStartHandler == nil {
		return errors.New("No start handler registered")
	}

	entryList, err := journal.Initialize()

	if err != nil {
		return err
	}

	entriesMu.Lock()
	entries = entryList
	entriesMu.Unlock()

	server.OnOperationCreate(create)
	server.OnOperationAbort(abort)
	server.OnOperationStatus(status)
	server.OnListOperations(list)

	for _, e := range getEntries() {
		err = registerdStartHandler(e.Operation)

		if err != nil {
			UpdateOperationStatus(e.Operation.ID, &op.OperationStatus{
				Finished: true,
				ExitCode: 255,
				Message:  err.Error(),
			})
			return err
		}
	}

	return nil
}
