package protocol

import (
	"barf/internal/op"
)

// Message defines the parameters for a message
type Message struct {
	RequestCreate  *RequestCreate  `json:"requestCreate,omitempty"`
	ResponseCreate *ResponseCreate `json:"responseCreate,omitempty"`
	RequestAbort   *RequestAbort   `json:"requestAbort,omitempty"`
	ResponseAbort  *ResponseAbort  `json:"responseAbort,omitempty"`
	RequestStatus  *RequestStatus  `json:"requestStatus,omitempty"`
	ResponseStatus *ResponseStatus `json:"responseStatus,omitempty"`
	RequestList    *RequestList    `json:"requestList,omitempty"`
	ResponseList   *ResponseList   `json:"responseList,omitempty"`
	EventStatus    *EventStatus    `json:"eventStatus,omitempty"`
	EventOperation *EventOperation `json:"eventOperation,omitempty"`
}

// Response result types
const (
	ResponseOk    = 0
	ResponseError = 1
)

// EventID is a unique identifier for the event
type EventID string

// RequestID is a unique identifier for the request
type RequestID string

// ResponseResult indicates the success or failure of a request
type ResponseResult int

// ResponseMessage is a string message to explain the response result
type ResponseMessage string

// RequestCreate defines the parameters for an operation create request
type RequestCreate struct {
	ID   RequestID        `json:"id"`
	Type op.OperationType `json:"type"`
	Args op.OperationArgs `json:"args"`
}

// ResponseCreate defines the parameters for an operation create response
type ResponseCreate struct {
	ID        RequestID       `json:"id"`
	Result    ResponseResult  `json:"result"`
	Message   ResponseMessage `json:"message,omitempty"`
	Operation *op.Operation   `json:"operation"`
}

// RequestAbort defines the parameters for an operation abort request
type RequestAbort struct {
	ID          RequestID      `json:"id"`
	OperationID op.OperationID `json:"operationId"`
}

// ResponseAbort defines the parameters for an operation abort response
type ResponseAbort struct {
	ID      RequestID       `json:"id"`
	Result  ResponseResult  `json:"result"`
	Message ResponseMessage `json:"message,omitempty"`
}

// RequestStatus defines the parameters for an operation status request
type RequestStatus struct {
	ID          RequestID      `json:"id"`
	OperationID op.OperationID `json:"operationId"`
}

// ResponseStatus defines the parameters for an operation status response
type ResponseStatus struct {
	ID      RequestID       `json:"id"`
	Result  ResponseResult  `json:"result"`
	Message ResponseMessage `json:"message,omitempty"`

	Status *op.OperationStatus `json:"status"`
}

// RequestList defines the parameters for an operation listing request
type RequestList struct {
	ID RequestID `json:"id"`
}

// ResponseList defines the parameters for an operation listing response
type ResponseList struct {
	ID      RequestID       `json:"id"`
	Result  ResponseResult  `json:"result"`
	Message ResponseMessage `json:"message,omitempty"`

	Operations []*op.Operation `json:"operations"`
}

// EventStatus defines the parameters for an operation status event
type EventStatus struct {
	ID EventID `json:"id"`

	OperationID op.OperationID      `json:"operationId"`
	Status      *op.OperationStatus `json:"status"`
}

// EventOperation defines the parameters for a new operation event
type EventOperation struct {
	ID EventID `json:"id"`

	Operation *op.Operation `json:"operation"`
}
