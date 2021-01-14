package server

import (
	"fmt"
	"os"

	"barf/internal/com/channel"
	"barf/internal/com/protocol"
	"barf/internal/com/socket"
)

func onRequestAbortError(socket *socket.Socket, requestID protocol.RequestID, resultMessage string) {
	message := protocol.Message{
		ResponseAbort: &protocol.ResponseAbort{
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

	message := protocol.Message{
		ResponseAbort: &protocol.ResponseAbort{
			ID:     requestAbort.ID,
			Result: protocol.ResponseOk,
		},
	}

	err = channel.Send(socket, &message)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
