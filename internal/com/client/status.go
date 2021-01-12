package client

import (
	"errors"
	"time"

	"rft/internal/com/channel"
	"rft/internal/com/protocol"
	"rft/internal/op"
)

// OperationStatus sends a status request
func OperationStatus(operationID op.OperationID) (*op.OperationStatus, error) {
	requestID := protocol.GenerateRequestID()
	message := protocol.Message{
		RequestStatus: &protocol.RequestStatus{
			ID:          requestID,
			OperationID: operationID,
		},
	}

	err := channel.Broadcast(&message)

	if err != nil {
		return nil, err
	}

	c1 := make(chan *protocol.ResponseStatus)

	onResponse := func(response *protocol.ResponseStatus) {
		c1 <- response
	}

	bus.SubscribeOnce(string(requestID), onResponse)

	select {
	case res := <-c1:
		if res.Result == protocol.ResponseOk {
			return res.Status, nil
		}

		return nil, errors.New(string(res.Message))
	case <-time.After(10 * time.Second):
		bus.Unsubscribe(string(requestID), onResponse)
		return nil, errors.New("Timeout")
	}
}
