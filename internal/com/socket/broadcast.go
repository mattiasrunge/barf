package socket

// Broadcast data to all sockets
func Broadcast(data string) error {
	for _, socket := range sockets {
		err := socket.Write(data)

		if err != nil {
			return err
		}
	}

	return nil
}
