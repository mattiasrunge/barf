package socket

import (
	"fmt"
	"net"

	"rft/internal/config"
	"rft/internal/proc/life"
)

var clientSocket *Socket = nil

// Connect to a backend server via socket
func Connect() error {
	connection, err := net.Dial("unix", config.SocketFile)

	if err != nil {
		return err
	}

	life.AddShutdownHook(func() {
		Close()
		fmt.Println("Client socket closed.")
	})

	clientSocket = newSocket(connection)

	wg.Add(1)

	bus.Publish("new", clientSocket)

	clientSocket.OnClose(func() {
		bus.Publish("close")
	})

	clientSocket.OnError(func(err error) {
		bus.Publish("error", err)
	})

	bus.Publish("new", clientSocket)

	return nil
}
