package client

import (
	"errors"
	"time"

	"barf/internal/com/channel"
	"barf/internal/com/protocol"
	"barf/internal/op"
)

// ListOperations sends a list request
func ListOperations() ([]*op.Operation, error) {
	message := protocol.NewRequestListMessage()
	err := channel.Broadcast(message)

	if err != nil {
		return nil, err
	}

	c1 := make(chan *protocol.ResponseList)

	onResponse := func(response *protocol.ResponseList) {
		c1 <- response
	}

	bus.SubscribeOnce(string(message.RequestList.ID), onResponse)

	select {
	case res := <-c1:
		if res.Result == protocol.ResponseOk {
			return res.Operations, nil
		}

		return nil, errors.New(string(res.Message))
	case <-time.After(10 * time.Second):
		bus.Unsubscribe(string(message.RequestList.ID), onResponse)
		return nil, errors.New("timeout")
	}
}
