package journal

import "rft/internal/op"

type startHandler func(*op.Operation) error
type abortHandler func(op.OperationID) error

var registerdStartHandler startHandler = nil
var registerdAbortHandler abortHandler = nil

// OnOperationStart registers a handler for an operation start request
func OnOperationStart(handler startHandler) {
	registerdStartHandler = handler
}

// OnOperationAbort registers a handler for an operation abort request
func OnOperationAbort(handler abortHandler) {
	registerdAbortHandler = handler
}
