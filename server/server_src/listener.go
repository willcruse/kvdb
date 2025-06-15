package serversrc

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

// Listener is responsible for responding to incoming messages
type Listener interface {
	Listen() error
}

type TCPListener struct {
	Address string
	Storage StorageBackend
}

func (t *TCPListener) Listen() error {
	ln, err := net.Listen("tcp", t.Address)
	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("tcp_listener: Error accepting connection. %v\n", err)
			continue
		}

		go t.handleNetConn(conn)

	}
}

func (t *TCPListener) handleNetConn(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		commandValue, err := reader.ReadByte()
		if err != nil {
			log.Printf("handler_net_conn: Failed to read command value from stream. %v\n", err)
			return
		}

		commandValueInt := int(commandValue)
		key, err := readString(reader)
		if err != nil {
			log.Printf("handler_net_conn: Error reading key from stream %v\n", err)
			return
		}

		fmt.Printf("Command Value: %d\n", commandValueInt)

		switch commandValueInt {
		case GET_COMMAND:
			// handler get command
			fmt.Printf("Fetching %s\n", key)
			res, err := t.Storage.Get(key)
			if err != nil {
				return
			}
			fmt.Printf("Fetched %s -> %s\n", key, res)

			break
		case SET_COMMAND:
			value, err := readString(reader)
			if err != nil {
				log.Printf("handler_net_conn: Error reading value from stream %v\n", err)
				return
			}
			fmt.Printf("Setting %s to %s\n", key, value)
			err = t.Storage.Set(key, value)
			if err != nil {
				return
			}
			fmt.Printf("Set %s -> %s\n", key, value)
			break
		case DELETE_COMMAND:
			// handle delete command
			fmt.Printf("Deleting %s\n", key)
			err = t.Storage.Delete(key)
			if err != nil {
				return
			}
			fmt.Printf("Delete %s\n", key)
			break
		default:
			// handle unknown command
		}
	}

}

func readString(reader *bufio.Reader) (string, error) {
	keyLength, err := reader.ReadByte()
	if err != nil {
		err = fmt.Errorf("handler_net_conn: Failed to read length from stream. %w", err)
		return "", err
	}

	keyBuf := make([]byte, int(keyLength))
	nRead, err := io.ReadFull(reader, keyBuf)
	if err != nil {
		err = fmt.Errorf("handler_net_conn: Failed to read value from stream. %w", err)
		return "", err
	}
	if nRead != int(keyLength) {
		err = fmt.Errorf("handler_net_conn: Incorrect number of bytes read from stream")
		return "", err
	}

	return string(keyBuf), nil
}
