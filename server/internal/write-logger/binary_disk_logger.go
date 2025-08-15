package writelogger

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/willcruse/kvdb/server/v2/internal/commands"
)

type BinaryDiskLogger struct {
	FileName string
	file     *os.File
}

func (bdl *BinaryDiskLogger) Init() error {
	if bdl.FileName == "" {
		return fmt.Errorf("(BinaryDiskLogger) BinaryDiskLogger has no FileName property")
	}
	file, err := os.OpenFile(bdl.FileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		err = fmt.Errorf("(BinaryDiskLogger) Failed to open %s. Error: %v", bdl.FileName, err)
		return err
	}

	bdl.file = file
	return nil
}

func (bdl *BinaryDiskLogger) Close() error {
	if bdl.file == nil {
		err := errors.New("(BinaryDiskLogger) No file to close. File is already nil")
		return err
	}

	err := bdl.file.Close()
	if err != nil {
		err = fmt.Errorf("(BinaryDiskLogger) Failed to close file. Error: %+v", err)
		return err
	}

	bdl.file = nil
	return nil
}

func (bdl *BinaryDiskLogger) LogSet(key, value string) error {
	// Convert command to binary format
	// COMMAND | LEN | KEY | LEN | VALUE
	encodedKey := encodeString(key)
	encodedValue := encodeString(value)

	var toLog []byte
	toLog = make([]byte, 1)
	toLog[0] = byte(commands.SET_COMMAND)

	toLog = append(toLog, encodedKey...)
	toLog = append(toLog, encodedValue...)
	return bdl.logToFile(toLog)
}

func (bdl *BinaryDiskLogger) LogDelete(key string) error {
	// Convert command to binary format
	// COMMAND | LEN | KEY | LEN | VALUE
	commandByte := byte(commands.DELETE_COMMAND)
	encodedKey := encodeString(key)

	var toLog []byte
	toLog = make([]byte, 1+len(encodedKey))
	toLog[0] = commandByte

	toLog = append(toLog, encodedKey...)
	return bdl.logToFile(toLog)
}

func (bdl *BinaryDiskLogger) Replay() ([]commands.Command, error) {
	if bdl.file == nil {
		return nil, errors.New("(BinaryDiskLogger) Could replay commands. File is nil")
	}

	reader := bufio.NewReader(bdl.file)
	replayed := []commands.Command{}
	for {
		commandByte, err := reader.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return replayed, fmt.Errorf("(BinaryDiskLogger) Failed to read command byte. Err: '%v'", err)
		}

		fmt.Println(commandByte)

		switch int(commandByte) {
		case commands.SET_COMMAND:
			key, err := decodeString(reader)
			if err != nil {
				return replayed, fmt.Errorf("(BinaryDiskLogger) Failed to read key in set command. Err: '%w'", err)
			}

			value, err := decodeString(reader)
			if err != nil {
				return replayed, fmt.Errorf("(BinaryDiskLogger) Failed to read value in set command. Err: '%w'", err)
			}

			replayed = append(replayed, commands.CreateSetCommand(key, value))
		case commands.DELETE_COMMAND:
			key, err := decodeString(reader)
			if err != nil {
				return replayed, fmt.Errorf("(BinaryDiskLogger) Failed to read key in delete command. Err: '%w'", err)
			}

			replayed = append(replayed, commands.CreateDeleteCommand(key))
		default:
			return nil, fmt.Errorf("(BinaryDiskLogger): Unknown command %d\n", int(commandByte))
		}
	}

	return replayed, nil
}

func decodeString(reader *bufio.Reader) (string, error) {
	length, err := reader.ReadByte()
	if err != nil {
		err = fmt.Errorf("(BinaryDiskLogger): Failed to read length from stream. %w", err)
		return "", err
	}

	buffer := make([]byte, int(length))
	nRead, err := io.ReadFull(reader, buffer)
	if err != nil {
		err = fmt.Errorf("(BinaryDiskLogger): Failed to read value from stream. %w", err)
		return "", err
	}
	if nRead != int(length) {
		err = fmt.Errorf("(BinaryDiskLogger): Incorrect number of bytes read from stream")
		return "", err
	}

	return string(buffer), nil

}

func encodeString(toEncode string) []byte {
	encoded := make([]byte, len(toEncode)+1)

	encoded[0] = byte(len(toEncode))
	offset := 1
	for idx, ch := range toEncode {
		encoded[offset+idx] = byte(ch)
	}

	return encoded
}

func (bdl *BinaryDiskLogger) logToFile(toLog []byte) error {
	_, err := bdl.file.Write(toLog)
	if err != nil {
		err = fmt.Errorf("(StringDiskLogger) Failed to save log to file. Error: %v", err)
		return err
	}
	return nil
}
