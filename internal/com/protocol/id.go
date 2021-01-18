package protocol

import (
	"github.com/rs/xid"
)

func generateEventID() EventID {
	id := xid.New()

	return EventID(id.String())
}

func generateRequestID() RequestID {
	id := xid.New()

	return RequestID(id.String())
}
