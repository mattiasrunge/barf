package channel

import (
	"barf/internal/com/protocol"
	"barf/internal/com/socket"
)

type messageHandler func(*socket.Socket, *protocol.Message)

// OnMessage registers handler for messages
func OnMessage(handler messageHandler) {
	bus.Subscribe("message", handler)
}
