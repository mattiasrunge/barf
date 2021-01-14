package channel

import (
	"barf/internal/com/protocol"
	"barf/internal/com/socket"
)

// Broadcast sends message to all sockets
func Broadcast(message *protocol.Message) error {
	data, err := protocol.Encode(message)

	if err != nil {
		return err
	}

	return socket.Broadcast(data)
}

// Send message to socket
func Send(s *socket.Socket, message *protocol.Message) error {
	// TODO assert message only contains event

	data, err := protocol.Encode(message)

	if err != nil {
		return err
	}

	return s.Write(data)
}
