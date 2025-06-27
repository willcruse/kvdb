package commands

import (
	"fmt"
)

const (
	GET_COMMAND    = 0
	SET_COMMAND    = 1
	DELETE_COMMAND = 2

	NO_ERROR_ERROR_CODE      = 0
	SERVER_ERROR_ERROR_CODE  = 1
	USER_ERROR_ERROR_CODE    = 2
	UNKNOWN_ERROR_ERROR_CODE = 3
)

type Command struct {
	Identifier int
	Key        string
	Value      string
}

func CreateGetCommand(key string) Command {
	return Command{GET_COMMAND, key, ""}
}

func CreateSetCommand(key, value string) Command {
	return Command{SET_COMMAND, key, value}
}

func CreateDeleteCommand(key string) Command {
	return Command{DELETE_COMMAND, key, ""}
}

type Response struct {
	ErrorCode uint8
	Message   string
}

func (r *Response) Encode() ([]byte, error) {
	// No message set only send error code
	if len(r.Message) == 0 {
		return []byte{byte(r.ErrorCode)}, nil
	}

	// Message can be max 255 bytes
	// Size must be a single byte
	// Note: This does not deal with UTF-8 only ASCII
	if len(r.Message) > 255 {
		return []byte{}, fmt.Errorf("Message too big. Max message size is 255 got message of length %d", len(r.Message))
	}

	var encoded []byte
	encoded = make([]byte, 2+len(r.Message))

	encoded[0] = byte(r.ErrorCode)
	encoded[1] = byte(len(r.Message))
	offset := 2
	for idx, ch := range r.Message {
		encoded[idx+offset] = byte(ch)
	}

	return encoded, nil
}
