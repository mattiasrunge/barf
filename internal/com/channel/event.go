package channel

import (
	"rft/internal/com/protocol"
	"rft/internal/com/socket"
)

type messageHandler func(*socket.Socket, *protocol.Message)

// OnMessage registers handler for messages
func OnMessage(handler messageHandler) {
	bus.Subscribe("message", handler)
}
