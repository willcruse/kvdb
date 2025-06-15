package src

import "errors"

const (
	NO_ERROR = 0
)

func DecodeResponse(res []byte) (*Response, error) {
	return &Response{}, nil
	resSize := len(res)

	if resSize == 0 {
		return nil, errors.New("decode_response: no bytes to decode!")
	}

	errorCode := int(res[0])

	var value string
	if resSize == 1 {
		value = ""
	} else {
		valueSize := int(res[1])
		var valueBytes []byte
		valueBytes = make([]byte, valueSize)
		numCopied := copy(valueBytes[0:], res[2:])
		if numCopied != valueSize {
			return nil, errors.New("decode_response: an unexpected number of bytes was returned.")
		}

		value = string(valueBytes)
	}

	return &Response{ErrorCode: errorCode, Value: value}, nil
}

type Response struct {
	ErrorCode int
	Value     string
}
