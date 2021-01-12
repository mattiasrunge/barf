package server

import (
	"rft/internal/com/channel"
	"rft/internal/com/protocol"
	"rft/internal/op"
)

// OperationCreated broadcasts a operation event
func OperationCreated(operation *op.Operation) error {
	messageID := protocol.GenerateEventID()
	message := protocol.Message{
		EventOperation: &protocol.EventOperation{
			ID:        messageID,
			Operation: operation,
		},
	}

	return channel.Broadcast(&message)
}

// OperationStatus broadcasts a status event
func OperationStatus(operationID op.OperationID, status *op.OperationStatus) error {
	messageID := protocol.GenerateEventID()
	message := protocol.Message{
		EventStatus: &protocol.EventStatus{
			ID:          messageID,
			OperationID: operationID,
			Status:      status,
		},
	}

	return channel.Broadcast(&message)
}
