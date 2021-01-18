package server

import (
	"barf/internal/com/channel"
	"barf/internal/com/protocol"
	"barf/internal/op"
)

// OperationCreated broadcasts a operation event
func OperationCreated(operation *op.Operation) error {
	message := protocol.NewEventOperationMessage(operation)

	return channel.Broadcast(message)
}

// OperationStatus broadcasts a status event
func OperationStatus(operationID op.OperationID, status *op.OperationStatus) error {
	message := protocol.NewEventStatusMessage(operationID, status)

	return channel.Broadcast(message)
}
