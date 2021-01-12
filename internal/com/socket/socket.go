package socket

import (
	"bufio"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/asaskevich/EventBus"
)

// Socket represents a data channel
type Socket struct {
	connection net.Conn
	bus        EventBus.Bus
	wg         sync.WaitGroup
}

func newSocket(c net.Conn) *Socket {
	socket := Socket{
		connection: c,
		bus:        EventBus.New(),
		wg:         sync.WaitGroup{},
	}

	socket.wg.Add(1)

	go func() {
		r := io.Reader(socket.connection)

		scanner := bufio.NewScanner(r)

		for scanner.Scan() {
			line := scanner.Text()
			// fmt.Println("RECV:", line)
			socket.bus.Publish("data", line)
		}

		err := scanner.Err()

		if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			socket.bus.Publish("error", err)
		}

		socket.bus.Publish("close")
		socket.wg.Done()
	}()

	registerSocket(&socket)

	return &socket
}

// Close the socket
func (s *Socket) Close() error {
	return s.connection.Close()
}

// Write to the socket
func (s *Socket) Write(data string) error {
	// fmt.Println("SEND:", data)
	_, err := s.connection.Write([]byte(data + "\n"))

	return err
}

// OnError registers a handler for error events
func (s *Socket) OnError(handler interface{}) error {
	return s.bus.Subscribe("error", handler)
}

// OnClose registers a handler for close events
func (s *Socket) OnClose(handler interface{}) error {
	return s.bus.Subscribe("close", handler)
}

// OnData registers a handler for data events
func (s *Socket) OnData(handler interface{}) error {
	return s.bus.Subscribe("data", handler)
}

func (s *Socket) waitOnClose() {
	s.wg.Wait()
}
