package channel

import (
	"fmt"
	"os"

	"barf/internal/com/protocol"
	"barf/internal/com/socket"

	"github.com/asaskevich/EventBus"
)

var bus EventBus.Bus

func onData(socket *socket.Socket, data string) {
	message, err := protocol.Decode(data)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	bus.Publish("message", socket, message)
}

func onSocket(socket *socket.Socket) {
	socket.OnData(func(data string) {
		onData(socket, data)
	})
}

func init() {
	bus = EventBus.New()

	socket.OnSocket(onSocket)
}
