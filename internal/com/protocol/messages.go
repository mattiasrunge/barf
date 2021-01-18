package protocol

import (
	"barf/internal/config"
	"barf/internal/op"
)

// Message defines the parameters for a message
type Message struct {
	Version        string          `json:"version"`
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

// NewRequestCreateMessage creates a new request create message
func NewRequestCreateMessage(opType op.OperationType, args op.OperationArgs) *Message {
	requestID := generateRequestID()
	message := Message{
		Version: config.Version,
		RequestCreate: &RequestCreate{
			ID:   requestID,
			Type: opType,
			Args: args,
		},
	}

	return &message
}

// NewRequestAbortMessage creates a new request abort message
func NewRequestAbortMessage(operationID op.OperationID) *Message {
	requestID := generateRequestID()
	message := Message{
		Version: config.Version,
		RequestAbort: &RequestAbort{
			ID:          requestID,
			OperationID: operationID,
		},
	}

	return &message
}

// NewRequestStatusMessage creates a new request status message
func NewRequestStatusMessage(operationID op.OperationID) *Message {
	requestID := generateRequestID()
	message := Message{
		Version: config.Version,
		RequestStatus: &RequestStatus{
			ID:          requestID,
			OperationID: operationID,
		},
	}

	return &message
}

// NewRequestListMessage creates a new request list message
func NewRequestListMessage() *Message {
	requestID := generateRequestID()
	message := Message{
		Version: config.Version,
		RequestList: &RequestList{
			ID: requestID,
		},
	}

	return &message
}

// NewResponseAbortMessage creates a new response abort message
func NewResponseAbortMessage(requestID RequestID, result ResponseResult, msg ResponseMessage) *Message {
	message := Message{
		Version: config.Version,
		ResponseAbort: &ResponseAbort{
			ID:      requestID,
			Result:  result,
			Message: msg,
		},
	}

	return &message
}

// NewResponseCreateMessage creates a new response create message
func NewResponseCreateMessage(requestID RequestID, result ResponseResult, msg ResponseMessage, operation *op.Operation) *Message {
	message := Message{
		Version: config.Version,
		ResponseCreate: &ResponseCreate{
			ID:        requestID,
			Result:    result,
			Message:   msg,
			Operation: operation,
		},
	}

	return &message
}

// NewResponseListMessage creates a new response list message
func NewResponseListMessage(requestID RequestID, result ResponseResult, msg ResponseMessage, operations []*op.Operation) *Message {
	message := Message{
		Version: config.Version,
		ResponseList: &ResponseList{
			ID:         requestID,
			Result:     result,
			Message:    msg,
			Operations: operations,
		},
	}

	return &message
}

// NewResponseStatusMessage creates a new response status message
func NewResponseStatusMessage(requestID RequestID, result ResponseResult, msg ResponseMessage, status *op.OperationStatus) *Message {
	message := Message{
		Version: config.Version,
		ResponseStatus: &ResponseStatus{
			ID:      requestID,
			Result:  result,
			Message: msg,
			Status:  status,
		},
	}

	return &message
}

// NewEventOperationMessage creates a new event operation message
func NewEventOperationMessage(operation *op.Operation) *Message {
	eventID := generateEventID()
	message := Message{
		Version: config.Version,
		EventOperation: &EventOperation{
			ID:        eventID,
			Operation: operation,
		},
	}

	return &message
}

// NewEventStatusMessage creates a new event status message
func NewEventStatusMessage(operationID op.OperationID, status *op.OperationStatus) *Message {
	eventID := generateEventID()
	message := Message{
		Version: config.Version,
		EventStatus: &EventStatus{
			ID:          eventID,
			OperationID: operationID,
			Status:      status,
		},
	}

	return &message
}
