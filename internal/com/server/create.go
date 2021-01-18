package server

import (
	"fmt"
	"os"

	"barf/internal/com/channel"
	"barf/internal/com/protocol"
	"barf/internal/com/socket"
)

func onRequestCreateError(socket *socket.Socket, requestID protocol.RequestID, resultMessage string) {
	message := protocol.NewResponseCreateMessage(requestID, protocol.ResponseError, protocol.ResponseMessage(resultMessage), nil)
	err := channel.Send(socket, message)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func onRequestCreate(socket *socket.Socket, requestCreate *protocol.RequestCreate) {
	if registerdCreateHandler == nil {
		onRequestCreateError(socket, requestCreate.ID, "No create handler registered")

		return
	}

	operation, err := registerdCreateHandler(requestCreate.Type, requestCreate.Args)

	if err != nil {
		onRequestCreateError(socket, requestCreate.ID, err.Error())

		return
	}

	message := protocol.NewResponseCreateMessage(requestCreate.ID, protocol.ResponseOk, "", operation)
	err = channel.Send(socket, message)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
