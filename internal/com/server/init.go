package server

import (
	"rft/internal/com/channel"
	"rft/internal/com/protocol"
	"rft/internal/com/socket"
)

func onMessage(socket *socket.Socket, message *protocol.Message) {
	if message.RequestCreate != nil {
		onRequestCreate(socket, message.RequestCreate)
	} else if message.RequestAbort != nil {
		onRequestAbort(socket, message.RequestAbort)
	} else if message.RequestStatus != nil {
		onRequestStatus(socket, message.RequestStatus)
	} else if message.RequestList != nil {
		onRequestList(socket, message.RequestList)
	}
}

func init() {
	channel.OnMessage(onMessage)
}
