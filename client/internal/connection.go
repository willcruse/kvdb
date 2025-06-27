package internal

import (
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
}

func CreateTCPServerConnection(serverAddress string) (*TCPServerConnection, error) {
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return nil, fmt.Errorf("tcp_conn: failed to create TCP connection. %w", err)
	}

	tcpServerConnection := TCPServerConnection{serverAddress, conn}
	return &tcpServerConnection, nil
}

func (c *TCPServerConnection) SendMessage(message []byte) ([]byte, error) {
	numWritten, err := c.socketConnection.Write(message)
	if err != nil {
		return nil, fmt.Errorf("tcp_conn: error writing message to connection. %w", err)
	}

	if numWritten != len(message) {
		return nil, errors.New("tcp_conn: failed to write entire message to connection")
	}

	var readBuffer []byte
	_, err = c.socketConnection.Read(readBuffer)

	if err != nil {
		return nil, fmt.Errorf("tcp_conn: failed to read from connection. %w", err)
	}

	// if numRead == 0 {
	// 	return nil, errors.New("tcp_conn: no bytes read from connection.")
	// }

	return readBuffer, nil
}
