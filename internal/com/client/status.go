package client

import (
	"errors"
	"time"

	"barf/internal/com/channel"
	"barf/internal/com/protocol"
	"barf/internal/op"
)

// OperationStatus sends a status request
func OperationStatus(operationID op.OperationID) (*op.OperationStatus, error) {
	message := protocol.NewRequestStatusMessage(operationID)
	err := channel.Broadcast(message)

	if err != nil {
		return nil, err
	}

	c1 := make(chan *protocol.ResponseStatus)

	onResponse := func(response *protocol.ResponseStatus) {
		c1 <- response
	}

	bus.SubscribeOnce(string(message.RequestStatus.ID), onResponse)

	select {
	case res := <-c1:
		if res.Result == protocol.ResponseOk {
			return res.Status, nil
		}

		return nil, errors.New(string(res.Message))
	case <-time.After(10 * time.Second):
		bus.Unsubscribe(string(message.RequestStatus.ID), onResponse)
		return nil, errors.New("timeout")
	}
}
