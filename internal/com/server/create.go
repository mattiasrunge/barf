package server

import (
	"fmt"
	"os"

	"barf/internal/com/channel"
	"barf/internal/com/protocol"
	"barf/internal/com/socket"
)

func onRequestCreateError(socket *socket.Socket, requestID protocol.RequestID, resultMessage string) {
	message := protocol.Message{
		ResponseCreate: &protocol.ResponseCreate{
			ID:      requestID,
			Result:  protocol.ResponseError,
			Message: protocol.ResponseMessage(resultMessage),
		},
	}

	err := channel.Send(socket, &message)

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

	message := protocol.Message{
		ResponseCreate: &protocol.ResponseCreate{
			ID:        requestCreate.ID,
			Result:    protocol.ResponseOk,
			Operation: operation,
		},
	}

	err = channel.Send(socket, &message)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
