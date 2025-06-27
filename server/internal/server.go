package internal

import (
	"bufio"
	"fmt"
	"io"
	"log"

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

	commandsToReplay, err := server.WriteLogger.Replay()
	if err != nil {
		err = fmt.Errorf("(Server) Failed to replay WriteLogger. Error: %w", err)
		return err
	}

	for _, commandToReplay := range commandsToReplay {
		log.Printf("Command Replay: %+v\n", commandToReplay)
		switch commandToReplay.Identifier {
		case commands.SET_COMMAND:
			err = server.StorageBackend.Set(commandToReplay.Key, commandToReplay.Value)
			if err != nil {
				err = fmt.Errorf("(Server): Failed to apply SET command. Key = %s Value = %s. Error: %w", commandToReplay.Key, commandToReplay.Value, err)
				return err
			}
		case commands.DELETE_COMMAND:
			err = server.StorageBackend.Delete(commandToReplay.Key)
			if err != nil {
				err = fmt.Errorf("(Server): Failed to apply DELETE command. Key = %s. Error: %w", commandToReplay.Key, err)
				return err
			}
		default:
			return fmt.Errorf("(Server) Unknown command to reply %+v", commandToReplay)
		}
	}

	return nil
}

func (server *Server) Listen() error {
	defer server.WriteLogger.Close()
	connChan := make(chan listener.Readable, CONNECTION_CHANNEL_BUFFER_SIZE)
	go server.messageReceiver(connChan)
	err := server.Listener.Listen(connChan)
	if err != nil {
		err = fmt.Errorf("(Server) Failed to listen. Error: %w", err)
		return err
	}

	return nil
}

func (server *Server) messageReceiver(connChan chan listener.Readable) {
	for conn := range connChan {
		// TODO: Handle concurrency issues
		go server.handleConnection(conn)
	}
}

func (server *Server) handleConnection(conn listener.Readable) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		response := commands.Response{ErrorCode: commands.NO_ERROR_ERROR_CODE, Message: ""}
		commandValue, err := reader.ReadByte()

		if err == io.EOF {
			log.Println("handler_net_conn: Connection closed")
			break
		}

		if err != nil {
			log.Printf("handler_net_conn: Failed to read command value from stream. %v\n", err)
			response.ErrorCode = commands.SERVER_ERROR_ERROR_CODE
			break
		}

		commandValueInt := int(commandValue)
		key, err := readString(reader)
		if err != nil {
			log.Printf("handler_net_conn: Error reading key from stream %v\n", err)
			response.ErrorCode = commands.SERVER_ERROR_ERROR_CODE
			break
		}

		fmt.Printf("Command Value: %d\n", commandValueInt)

		switch commandValueInt {
		case commands.GET_COMMAND:
			fmt.Printf("Fetching %s\n", key)
			res, err := server.StorageBackend.Get(key)
			if err != nil {
				log.Printf("handler_net_conn: Error fetching from storage backend %v\n", err)
				response.ErrorCode = commands.SERVER_ERROR_ERROR_CODE
				break
			}
			fmt.Printf("Fetched %s -> %s\n", key, res)
			response.Message = res

		case commands.SET_COMMAND:
			value, err := readString(reader)
			if err != nil {
				log.Printf("handler_net_conn: Error reading value from stream %v\n", err)
				response.ErrorCode = commands.SERVER_ERROR_ERROR_CODE
				break
			}
			fmt.Printf("Setting %s to %s\n", key, value)
			err = server.StorageBackend.Set(key, value)
			if err != nil {
				log.Printf("handler_net_conn: Error setting value %v\n", err)
				response.ErrorCode = commands.SERVER_ERROR_ERROR_CODE
				break
			}

			err = server.WriteLogger.LogSet(key, value)
			if err != nil {
				log.Printf("handler_net_conn: Error logging operation %v\n", err)
				response.ErrorCode = commands.SERVER_ERROR_ERROR_CODE
				break
			}
			fmt.Printf("Set %s -> %s\n", key, value)

		case commands.DELETE_COMMAND:
			fmt.Printf("Deleting %s\n", key)
			err = server.StorageBackend.Delete(key)
			if err != nil {
				log.Printf("handler_net_conn: Error deleting value %v\n", err)
				response.ErrorCode = commands.SERVER_ERROR_ERROR_CODE
				break
			}
			err = server.WriteLogger.LogDelete(key)
			if err != nil {
				log.Printf("handler_net_conn: Error logging operation %v\n", err)
				response.ErrorCode = commands.SERVER_ERROR_ERROR_CODE
				break
			}
			fmt.Printf("Delete %s\n", key)

		default:
			response.ErrorCode = commands.USER_ERROR_ERROR_CODE
		}

		err = sendResponse(response, conn)
		if err != nil {
			log.Printf("handler_net_conn: Error sending response %v\n", err)
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

func sendResponse(response commands.Response, conn listener.Readable) error {
	encoded, err := response.Encode()
	if err != nil {
		err = fmt.Errorf("(Server) Failed to encode response. Error: %w", err)
		return err
	}

	log.Println("Sending response")
	_, err = conn.Write(encoded)
	if err != nil {
		err = fmt.Errorf("(Server) Failed to send response. Error: %w", err)
		return err
	}

	return nil
}
