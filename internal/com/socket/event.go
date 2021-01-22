package socket

type errorHandler func(error)
type closeHandler func(bool)
type socketHandler func(*Socket)

// OnError registers handler for error events
func OnError(handler errorHandler) {
	bus.Subscribe("error", handler)
}

// OnClose registers handler for close event
func OnClose(handler closeHandler) {
	bus.Subscribe("close", handler)
}

// OnSocket registers handler for new socket connections
func OnSocket(handler socketHandler) {
	bus.Subscribe("new", handler)
}
