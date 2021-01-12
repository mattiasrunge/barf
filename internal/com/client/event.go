package client

import "rft/internal/op"

type operationCreatedHandler func(op.Operation)
type operationStatusHandler func(op.OperationID, *op.OperationStatus)

// OnOperationCreated registers handler for operation created events
func OnOperationCreated(handler operationCreatedHandler) {
	bus.Subscribe("operation", handler)
}

// OnOperationStatus registers handler for operation status events
func OnOperationStatus(handler operationStatusHandler) {
	bus.Subscribe("status", handler)
}
