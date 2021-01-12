package protocol

import "encoding/json"

// Encode the supplied message into bytes
func Encode(message *Message) (string, error) {
	data, err := json.Marshal(message)

	return string(data), err
}
