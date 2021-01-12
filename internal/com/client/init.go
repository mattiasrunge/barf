package client

import (
	"rft/internal/com/channel"
	"rft/internal/com/protocol"
	"rft/internal/com/socket"

	"github.com/asaskevich/EventBus"
)

var bus EventBus.Bus

func onMessage(socket *socket.Socket, message *protocol.Message) {
	if message.ResponseCreate != nil {
		bus.Publish(string(message.ResponseCreate.ID), message.ResponseCreate)
	} else if message.ResponseAbort != nil {
		bus.Publish(string(message.ResponseAbort.ID), message.ResponseAbort)
	} else if message.ResponseStatus != nil {
		bus.Publish(string(message.ResponseStatus.ID), message.ResponseStatus)
	} else if message.ResponseList != nil {
		bus.Publish(string(message.ResponseList.ID), message.ResponseList)
	} else if message.EventStatus != nil {
		bus.Publish("status", message.EventStatus.OperationID, message.EventStatus.Status)
	} else if message.EventOperation != nil {
		bus.Publish("operation", &message.EventOperation.Operation)
	}
}

func init() {
	bus = EventBus.New()

	channel.OnMessage(onMessage)
}
