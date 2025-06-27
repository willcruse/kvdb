package internal

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/willcruse/kvdb/server/v2/internal/commands"
	listener "github.com/willcruse/kvdb/server/v2/internal/listener"
	storagebackend "github.com/willcruse/kvdb/server/v2/internal/storage-backend"
	writelogger "github.com/willcruse/kvdb/server/v2/internal/write-logger"
)

const (
	CONNECTION_CHANNEL_BUFFER_SIZE = 128
)

type Server struct {
	Listener       listener.Listener
	StorageBackend storagebackend.StorageBackend
	WriteLogger    writelogger.WriteOperationLogger
}

func (server *Server) Init() error {
	server.StorageBackend.Init()
	err := server.WriteLogger.Init()
	if err != nil {
		err = fmt.Errorf("(Server) Failed to init WriteLogger. Error: %w", err)
		return err
	}

	return nil
}

func (server *Server) Listen() error {
	defer server.WriteLogger.Close()
	connChan := make(chan net.Conn, CONNECTION_CHANNEL_BUFFER_SIZE)
	go server.messageReceiver(connChan)
	err := server.Listener.Listen(connChan)
	if err != nil {
		err = fmt.Errorf("(Server) Failed to listen. Error: %w", err)
		return err
	}

	return nil
}

func (server *Server) messageReceiver(connChan chan net.Conn) {
	for conn := range connChan {
		go server.handleMessage(conn)
	}
}

func (server *Server) handleMessage(conn net.Conn) {
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
		case commands.GET_COMMAND:
			// handler get command
			fmt.Printf("Fetching %s\n", key)
			res, err := server.StorageBackend.Get(key)
			if err != nil {
				return
			}
			fmt.Printf("Fetched %s -> %s\n", key, res)

		case commands.SET_COMMAND:
			value, err := readString(reader)
			if err != nil {
				log.Printf("handler_net_conn: Error reading value from stream %v\n", err)
				return
			}
			fmt.Printf("Setting %s to %s\n", key, value)
			err = server.StorageBackend.Set(key, value)
			if err != nil {
				log.Printf("handler_net_conn: Error setting value %v\n", err)
				return
			}

			err = server.WriteLogger.LogSet(key, value)
			if err != nil {
				log.Printf("handler_net_conn: Error logging operation %v\n", err)
				return
			}
			fmt.Printf("Set %s -> %s\n", key, value)

		case commands.DELETE_COMMAND:
			fmt.Printf("Deleting %s\n", key)
			err = server.StorageBackend.Delete(key)
			if err != nil {
				log.Printf("handler_net_conn: Error deleting value %v\n", err)
				return
			}
			err = server.WriteLogger.LogDelete(key)
			if err != nil {
				log.Printf("handler_net_conn: Error logging operation %v\n", err)
				return
			}
			fmt.Printf("Delete %s\n", key)

		default:
			// TODO: handle unknown command
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
