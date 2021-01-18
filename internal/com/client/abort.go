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
	message := protocol.NewRequestAbortMessage(operationID)
	err := channel.Broadcast(message)

	if err != nil {
		return err
	}

	c1 := make(chan *protocol.ResponseAbort)

	onResponse := func(response *protocol.ResponseAbort) {
		c1 <- response
	}

	bus.SubscribeOnce(string(message.RequestAbort.ID), onResponse)

	select {
	case res := <-c1:
		if res.Result == protocol.ResponseOk {
			return nil
		}

		return errors.New(string(res.Message))
	case <-time.After(10 * time.Second):
		bus.Unsubscribe(string(message.RequestAbort.ID), onResponse)
		return errors.New("timeout")
	}
}
