package socket

import (
	"fmt"
	"net"
	"os"

	"rft/internal/config"
	"rft/internal/proc/life"
)

var serverSocket net.Listener = nil

// Listen starts listening for connections from clients
func Listen() error {
	err := os.RemoveAll(config.SocketFile)

	if err != nil {
		return err
	}

	serverSocket, err := net.Listen("unix", config.SocketFile)

	if err != nil {
		return err
	}

	wg.Add(1)

	life.AddShutdownHook(func() {
		Close()
		os.RemoveAll(config.SocketFile)
		fmt.Println("Server socket closed.")
	})

	go func() {
		for {
			connection, err := serverSocket.Accept()

			if err != nil {
				bus.Publish("error", err)
				bus.Publish("close")
			}

			fmt.Println("New client connected")
			socket := newSocket(connection)

			socket.bus.Subscribe("close", func() {
				fmt.Println("Client disconnected")
			})

			bus.Publish("new", socket)
		}
	}()

	return nil
}
