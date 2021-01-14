package client

import (
	"errors"
	"time"

	"barf/internal/com/channel"
	"barf/internal/com/protocol"
	"barf/internal/op"
)

// AbortOperation sends a abort request
func AbortOperation(operationID op.OperationID) error {
	requestID := protocol.GenerateRequestID()
	message := protocol.Message{
		RequestAbort: &protocol.RequestAbort{
			ID:          requestID,
			OperationID: operationID,
		},
	}

	err := channel.Broadcast(&message)

	if err != nil {
		return err
	}

	c1 := make(chan *protocol.ResponseAbort)

	onResponse := func(response *protocol.ResponseAbort) {
		c1 <- response
	}

	bus.SubscribeOnce(string(requestID), onResponse)

	select {
	case res := <-c1:
		if res.Result == protocol.ResponseOk {
			return nil
		}

		return errors.New(string(res.Message))
	case <-time.After(10 * time.Second):
		bus.Unsubscribe(string(requestID), onResponse)
		return errors.New("Timeout")
	}
}
