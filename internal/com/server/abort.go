package server

import (
	"fmt"
	"os"

	"barf/internal/com/channel"
	"barf/internal/com/protocol"
	"barf/internal/com/socket"
)

func onRequestAbortError(socket *socket.Socket, requestID protocol.RequestID, resultMessage string) {
	message := protocol.NewResponseAbortMessage(requestID, protocol.ResponseError, protocol.ResponseMessage(resultMessage))
	err := channel.Send(socket, message)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func onRequestAbort(socket *socket.Socket, requestAbort *protocol.RequestAbort) {
	if registerdAbortHandler == nil {
		onRequestAbortError(socket, requestAbort.ID, "No abort handler registered")

		return
	}

	err := registerdAbortHandler(requestAbort.OperationID)

	if err != nil {
		onRequestAbortError(socket, requestAbort.ID, err.Error())

		return
	}

	message := protocol.NewResponseAbortMessage(requestAbort.ID, protocol.ResponseOk, "")

	err = channel.Send(socket, message)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
