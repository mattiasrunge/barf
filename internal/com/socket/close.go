package socket

// Close the client or server socket
func Close() error {
	var err error = nil

	if clientSocket != nil {
		err = clientSocket.Close()
		clientSocket = nil
	} else if serverSocket != nil {
		for _, s := range sockets {
			s.Close()
		}

		err = serverSocket.Close()
		serverSocket = nil
	}

	wg.Done()

	return err
}

// WaitOnClose blocks until socket is closed
func WaitOnClose() {
	wg.Wait()
}
