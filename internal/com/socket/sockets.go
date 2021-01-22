package socket

import (
	"fmt"
)

var sockets []*Socket

func registerSocket(socket *Socket) {
	sockets = append(sockets, socket)

	socket.OnError(func(err error) {
		fmt.Println("Socket encountered error", err)
	})

	socket.OnClose(func(_ bool) {
		unregisterSocket(socket)
	})
}

func unregisterSocket(socket *Socket) {
	for i, s := range sockets {
		if s == socket {
			copy(sockets[i:], sockets[i+1:])
			sockets = sockets[:len(sockets)-1]
			return
		}
	}
}
