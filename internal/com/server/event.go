package server

import "barf/internal/op"

type createHandler func(op.OperationType, op.OperationArgs) (*op.Operation, error)
type abortHandler func(op.OperationID) error
type statusHandler func(op.OperationID) (*op.OperationStatus, error)
type listHandler func() ([]*op.Operation, error)

var registerdCreateHandler createHandler = nil
var registerdAbortHandler abortHandler = nil
var registerdStatusHandler statusHandler = nil
var registerdListHandler listHandler = nil

// OnOperationCreate registers a handler for operation create requests
func OnOperationCreate(handler createHandler) {
	registerdCreateHandler = handler
}

// OnOperationAbort registers a handler for operation abort requests
func OnOperationAbort(handler abortHandler) {
	registerdAbortHandler = handler
}

// OnOperationStatus registers a handler for operation status requests
func OnOperationStatus(handler statusHandler) {
	registerdStatusHandler = handler
}

// OnListOperations registers a handler for operation list requests
func OnListOperations(handler listHandler) {
	registerdListHandler = handler
}
