package protocol

import "encoding/json"

// Decode the supplied bytes into a message
func Decode(data string) (*Message, error) {
	var message Message

	err := json.Unmarshal([]byte(data), &message)

	return &message, err
}
