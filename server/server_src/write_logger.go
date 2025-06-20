package serversrc

import (
	"fmt"
	"os"
)

// TODO: Allow configuration
const (
	BUFFER_SIZE = 1024
	FILE_NAME   = "kvdb.db"
)

type WriteOperationLogger interface {
	Init() error
	LogSet(key, value string) error
	LogDelete(key string) error
}

// TODO: This should be replaced with a logger that stores actual commands in binary format
type StringDiskLogger struct {
	file *os.File
}

func (sdl *StringDiskLogger) Init() error {
	file, err := os.OpenFile(FILE_NAME, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		err = fmt.Errorf("(StringDiskLogger) Failed to open %s. Error: %v", FILE_NAME, err)
		return err
	}

	sdl.file = file
	return nil
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
