package server

import (
	"fmt"
	"os"

	"barf/internal/com/channel"
	"barf/internal/com/protocol"
	"barf/internal/com/socket"
)

func onRequestListError(socket *socket.Socket, requestID protocol.RequestID, resultMessage string) {
	message := protocol.Message{
		ResponseList: &protocol.ResponseList{
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

func onRequestList(socket *socket.Socket, requestList *protocol.RequestList) {
	if registerdListHandler == nil {
		onRequestListError(socket, requestList.ID, "No list handler registered")

		return
	}

	operations, err := registerdListHandler()

	if err != nil {
		onRequestListError(socket, requestList.ID, err.Error())

		return
	}

	message := protocol.Message{
		ResponseList: &protocol.ResponseList{
			ID:         requestList.ID,
			Result:     protocol.ResponseOk,
			Operations: operations,
		},
	}

	err = channel.Send(socket, &message)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
