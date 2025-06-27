package internal

import (
	"bufio"
	"errors"
	"fmt"
	"net"
)

type ServerConnection interface {
	SendMessage([]byte) ([]byte, error)
}

type TCPServerConnection struct {
	serverAddress    string
	socketConnection net.Conn
	reader           *bufio.Reader
	writer           *bufio.Writer
}

func CreateTCPServerConnection(serverAddress string) (*TCPServerConnection, error) {
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return nil, fmt.Errorf("tcp_conn: failed to create TCP connection. %w", err)
	}

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	tcpServerConnection := TCPServerConnection{
		serverAddress:    serverAddress,
		socketConnection: conn,
		reader:           reader,
		writer:           writer,
	}
	return &tcpServerConnection, nil
}

func (c *TCPServerConnection) SendMessage(message []byte) ([]byte, error) {
	numWritten, err := c.writer.Write(message)
	if err != nil {
		return nil, fmt.Errorf("tcp_conn: error writing message to connection. %w", err)
	}
	err = c.writer.Flush()
	if err != nil {
		return nil, fmt.Errorf("tcp_conn: error writing message to connection. %w", err)
	}

	if numWritten != len(message) {
		return nil, errors.New("tcp_conn: failed to write entire message to connection")
	}

	var readBuffer []byte
	// Allocating way too much memory
	// Max message size in theory is 257 bytes
	readBuffer = make([]byte, 512)
	// _, err = c.reader.Read(readBuffer)
	nRead, err := c.socketConnection.Read(readBuffer)

	if nRead == 0 {
		return nil, errors.New("tcp_conn: no bytes read")
	}

	if err != nil {
		return nil, fmt.Errorf("tcp_conn: failed to read from connection. %w", err)
	}

	return readBuffer, nil
}
