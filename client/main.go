package main

import (
	"fmt"
	"log"

	"github.com/willcruse/kvdb/client/v2/internal"
)

const SERVER_ADDRESS = "localhost:1337"

func main() {
	fmt.Printf("Connecting to server on %s\n", SERVER_ADDRESS)
	tcpConn, err := internal.CreateTCPServerConnection(SERVER_ADDRESS)
	if err != nil {
		log.Fatalf("Failed to start TCP Server Connection. Error: %v\n", err)
	}

	key := "test"
	value := "hello,world"
	setCommand := internal.SetCommand{Key: key, Value: value}
	encodedSetCommand, err := setCommand.Encode()
	if err != nil {
		log.Fatalf("Failed to encode set command. Error: %v\n", err)
	}

	getCommand := internal.GetCommand{Key: key}
	encodedGetCommand, err := getCommand.Encode()
	if err != nil {
		log.Fatalf("Failed to encode get command. Error: %v\n", err)
	}

	deleteCommand := internal.DeleteCommand{Key: key}
	encodedDeleteCommand, err := deleteCommand.Encode()
	if err != nil {
		log.Fatalf("Failed to encode delete command. Error: %v\n", err)
	}

	setRes, err := tcpConn.SendMessage(encodedSetCommand)
	if err != nil {
		log.Fatalf("Failed to send set command. Error: %v\n", err)
	}
	decodedSetRes, err := internal.DecodeResponse(setRes)
	if err != nil {
		log.Fatalf("Failed to decode set response. Error: %v\n", err)
	}
	fmt.Printf("Set Response: %+v\n", decodedSetRes)

	getRes, err := tcpConn.SendMessage(encodedGetCommand)
	if err != nil {
		log.Fatalf("Failed to send get command. Error: %v\n", err)
	}
	decodedGetRes, err := internal.DecodeResponse(getRes)
	if err != nil {
		log.Fatalf("Failed to decode get response. Error: %v\n", err)
	}
	fmt.Printf("Get Response: %+v\n", decodedGetRes)

	deleteRes, err := tcpConn.SendMessage(encodedDeleteCommand)
	if err != nil {
		log.Fatalf("Failed to send delete command. Error: %v\n", err)
	}
	decodedDeleteRes, err := internal.DecodeResponse(deleteRes)
	if err != nil {
		log.Fatalf("Failed to decode delete response. Error: %v\n", err)
	}
	fmt.Printf("Delete Response: %+v\n", decodedDeleteRes)
}
