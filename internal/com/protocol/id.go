package protocol

import (
	"github.com/rs/xid"
)

// GenerateEventID creates a unique event id
func GenerateEventID() EventID {
	id := xid.New()

	return EventID(id.String())
}

// GenerateRequestID creates a unique request id
func GenerateRequestID() RequestID {
	id := xid.New()

	return RequestID(id.String())
}
