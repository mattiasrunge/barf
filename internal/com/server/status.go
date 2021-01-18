package server

import (
	"fmt"
	"os"

	"barf/internal/com/channel"
	"barf/internal/com/protocol"
	"barf/internal/com/socket"
)

func onRequestStatusError(socket *socket.Socket, requestID protocol.RequestID, resultMessage string) {
	message := protocol.NewResponseStatusMessage(requestID, protocol.ResponseError, protocol.ResponseMessage(resultMessage), nil)
	err := channel.Send(socket, message)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func onRequestStatus(socket *socket.Socket, requestStatus *protocol.RequestStatus) {
	if registerdStatusHandler == nil {
		onRequestStatusError(socket, requestStatus.ID, "No status handler registered")

		return
	}

	status, err := registerdStatusHandler(requestStatus.OperationID)

	if err != nil {
		onRequestStatusError(socket, requestStatus.ID, err.Error())

		return
	}

	message := protocol.NewResponseStatusMessage(requestStatus.ID, protocol.ResponseOk, "", status)
	err = channel.Send(socket, message)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
