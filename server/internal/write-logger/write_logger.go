package writelogger

import (
	"bufio"
	"fmt"
	"os"

	"github.com/willcruse/kvdb/server/v2/internal/commands"
)

type WriteOperationLogger interface {
	Init() error
	Close() error
	LogSet(key, value string) error
	LogDelete(key string) error
	Replay() ([]commands.Command, error)
}

// TODO: This should be replaced with a logger that stores actual commands in binary format
type StringDiskLogger struct {
	FileName string
	file     *os.File
}

func (sdl *StringDiskLogger) Init() error {
	if sdl.FileName == "" {
		return fmt.Errorf("(StringDiskLogger) StringDiskLogger has no FileName property")
	}
	file, err := os.OpenFile(sdl.FileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		err = fmt.Errorf("(StringDiskLogger) Failed to open %s. Error: %v", sdl.FileName, err)
		return err
	}

	sdl.file = file
	return nil
}

func (sdl *StringDiskLogger) Close() error {
	err := sdl.file.Close()
	if err != nil {
		err = fmt.Errorf("(StringDiskLogger) Failed to close file. Error: %+v", err)
		return err
	}

	sdl.file = nil
	return nil
}

func (sdl *StringDiskLogger) Replay() ([]commands.Command, error) {
	if sdl.file == nil {
		err := fmt.Errorf("(StringDiskLogger) Call to Replay with no file set. Have you called `Init()`?")
		return nil, err
	}

	var logged_commands []commands.Command
	scanner := bufio.NewScanner(sdl.file)
	for scanner.Scan() {
		line := scanner.Text()
		switch line {
		case "SET":
			canScan := scanner.Scan()
			if !canScan {
				return nil, fmt.Errorf("(StringDiskLogger) Unexpected EOF")
			}
			keyLine := scanner.Text()
			key, err := sdl.parseLogLine(keyLine)
			if err != nil {
				return nil, err
			}
			canScan = scanner.Scan()
			if !canScan {
				return nil, fmt.Errorf("(StringDiskLogger) Unexpected EOF")
			}
			valueLine := scanner.Text()
			value, err := sdl.parseLogLine(valueLine)
			if err != nil {
				return nil, err
			}
			logged_commands = append(logged_commands, commands.CreateSetCommand(key, value))
		case "DELETE":
			canScan := scanner.Scan()
			if !canScan {
				return nil, fmt.Errorf("(StringDiskLogger) Unexpected EOF")
			}
			keyLine := scanner.Text()
			key, err := sdl.parseLogLine(keyLine)
			if err != nil {
				return nil, err
			}
			logged_commands = append(logged_commands, commands.CreateDeleteCommand(key))
		default:
			err := fmt.Errorf("(StringDiskLogger) Unexpected line content. Expected 'SET' or 'DELETE' got '%s'", line)
			return nil, err
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("(StringDiskLogger) Failed to read from file. Error: %w", err)
	}

	return logged_commands, nil
}

func (sdl *StringDiskLogger) parseLogLine(line string) (string, error) {
	parsingLength := true
	length := 0
	var parsed []rune
	for _, val := range line {
		if parsingLength {
			// We've reached the separator
			if val == ' ' {
				parsingLength = false
			} else {
				intVal := int(val - '0')
				if intVal < 0 || intVal > 9 {
					return "", fmt.Errorf("(StringDiskLogger) Failure to parse length of log line. Expected int got %d", val)
				}

				length = (length * 10) + intVal
			}
		} else {
			// TODO: Length validation logic
			parsed = append(parsed, val)
		}
	}

	if len(parsed) == 0 {
		return "", fmt.Errorf("(StringDiskLogger) Expected value from line got nothing. Line: %s", line)
	}

	return string(parsed), nil
}

func (sdl *StringDiskLogger) LogSet(key, value string) error {
	// String escaping is fun...
	// TODO: No string escaping so weird log format
	// Great for debugging but not ideal actual
	keyLen := len(key)
	valueLen := len(value)
	logStr := fmt.Sprintf("SET\n%d %s\n%d %s\n", keyLen, key, valueLen, value)
	return sdl.logToFile(logStr)
}

func (sdl *StringDiskLogger) LogDelete(key string) error {
	// String escaping is fun...
	// TODO: No string escaping so weird log format
	// Great for debugging but not ideal actual
	keyLen := len(key)
	logStr := fmt.Sprintf("DELETE\n%d %s\n", keyLen, key)
	return sdl.logToFile(logStr)
}

func (sdl *StringDiskLogger) logToFile(toLog string) error {
	_, err := sdl.file.WriteString(toLog)
	if err != nil {
		err = fmt.Errorf("(StringDiskLogger) Failed to save log to file. Error: %v", err)
		return err
	}
	return nil
}
