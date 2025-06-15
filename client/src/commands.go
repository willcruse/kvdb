package src

import (
	"errors"
)

const MAX_STRING_SIZE = 255 // uint8 MAX

const (
	GET_COMMAND    = 0
	SET_COMMAND    = 1
	DELETE_COMMAND = 2
)

type Command interface {
	Encode() ([]byte, error)
}

func encodeString(toEncode string) (int, []byte) {
	encoded := []byte(toEncode)
	return len(encoded), encoded
}

type GetCommand struct {
	Key string
}

func (g *GetCommand) Encode() ([]byte, error) {
	keySize, encodedKey := encodeString(g.Key)

	if keySize > MAX_STRING_SIZE {
		return nil, errors.New("get_command: Key size is greater then max size allowed (255)")
	}

	var encodedMessage []byte
	encodedMessage = make([]byte, len(encodedKey)+2) // 2 bytes used for Command type + length
	encodedMessage[0] = GET_COMMAND
	encodedMessage[1] = byte(keySize)
	numCopied := copy(encodedMessage[2:], encodedKey)

	if numCopied != keySize {
		return nil, errors.New("get_command: failed to copy full key into encoded message")
	}

	return encodedMessage, nil
}

type SetCommand struct {
	Key   string
	Value string
}

func (s *SetCommand) Encode() ([]byte, error) {
	keySize, encodedKey := encodeString(s.Key)
	if keySize > MAX_STRING_SIZE {
		return nil, errors.New("set_command: Key size is greater then max size allowed (255)")
	}

	valueSize, encodedValue := encodeString(s.Value)
	if valueSize > MAX_STRING_SIZE {
		return nil, errors.New("set_command: Value size is greater then max size allowed (255)")
	}

	totalMessageSize := 1 + 1 + keySize + 1 + valueSize

	var encodedMessage []byte
	encodedMessage = make([]byte, totalMessageSize) // 2 bytes used for Command type + length
	encodedMessage[0] = SET_COMMAND
	encodedMessage[1] = byte(keySize)
	numCopied := copy(encodedMessage[2:], encodedKey)
	if numCopied != keySize {
		return nil, errors.New("set_command: failed to copy full key into encoded message")
	}

	startValueIndex := 1 + 1 + keySize
	encodedMessage[startValueIndex] = byte(valueSize)
	numCopied = copy(encodedMessage[startValueIndex+1:], encodedValue)
	if numCopied != valueSize {
		return nil, errors.New("set_command: failed to copy full value into encoded message")
	}

	return encodedMessage, nil
}

type DeleteCommand struct {
	Key string
}

func (d *DeleteCommand) Encode() ([]byte, error) {
	keySize, encodedKey := encodeString(d.Key)
	if keySize > MAX_STRING_SIZE {
		return nil, errors.New("delete_command: Key size is greater then max size allowed (255)")
	}

	totalMessageSize := 1 + 1 + keySize

	var encodedMessage []byte
	encodedMessage = make([]byte, totalMessageSize) // 2 bytes used for Command type + length
	encodedMessage[0] = DELETE_COMMAND
	encodedMessage[1] = byte(keySize)
	numCopied := copy(encodedMessage[2:], encodedKey)
	if numCopied != keySize {
		return nil, errors.New("delete_command: failed to copy full key into encoded message")
	}

	return encodedMessage, nil
}
