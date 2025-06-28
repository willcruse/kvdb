package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/willcruse/kvdb/client/v2/internal"
)

const SERVER_ADDRESS = "localhost:1337"

const HELP_MESSAGE = `Commands
	GET <KEY>: Fetch value of <KEY> from server
	SET <KEY> <VALUE>: Set <KEY> to <VALUE>
	DELETE <KEY>: Delete <KEY> from server
	HELP: Print this message

	Note: Commands are case insensitive
`

func main() {
	fmt.Printf("Connecting to server on %s\n", SERVER_ADDRESS)
	tcpConn, err := internal.CreateTCPServerConnection(SERVER_ADDRESS)
	if err != nil {
		log.Fatalf("Failed to start TCP Server Connection. Error: %v\n", err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		line := scanner.Text()
		splitLine := strings.Split(line, " ")
		var command string
		if len(splitLine) == 0 {
			command = "help"
		} else {
			command = splitLine[0]
		}

		command = strings.ToUpper(command)

		switch command {
		case "GET":
			if len(splitLine) != 2 {
				fmt.Printf("GET command takes exactly one argument. Got %d.\n", len(splitLine)-1)
				continue
			}

			getCommand := internal.GetCommand{Key: splitLine[1]}
			encoded, err := getCommand.Encode()
			if err != nil {
				fmt.Printf("ERROR: Failed to encode GET command. Command: '%+v'. Error: %v\n", getCommand, err)
				continue
			}

			res, err := tcpConn.SendMessage(encoded)
			if err != nil {
				fmt.Printf("ERROR: Failed to send GET command. Command: '%+v'. Error: %v\n", getCommand, err)
				continue
			}

			decoded, err := internal.DecodeResponse(res)
			if err != nil {
				fmt.Printf("ERROR: Failed to decode response. Error: %v\n", err)
				continue
			}

			if decoded.ErrorCode != internal.NO_ERROR {
				fmt.Printf("ERROR: Server responeded with Error code %d\n", decoded.ErrorCode)
				continue
			}

			fmt.Printf("%s\n", decoded.Value)

		case "SET":
			if len(splitLine) != 3 {
				fmt.Printf("SET command takes two arguments. Got %d.\n", len(splitLine)-1)
				continue
			}

			setCommand := internal.SetCommand{Key: splitLine[1], Value: splitLine[2]}
			encoded, err := setCommand.Encode()
			if err != nil {
				fmt.Printf("ERROR: Failed to encode SET command. Command: '%+v'. Error: %v\n", setCommand, err)
				continue
			}

			res, err := tcpConn.SendMessage(encoded)
			if err != nil {
				fmt.Printf("ERROR: Failed to send SET command. Command: '%+v'. Error: %v\n", setCommand, err)
				continue
			}

			decoded, err := internal.DecodeResponse(res)
			if err != nil {
				fmt.Printf("ERROR: Failed to decode response. Error: %v\n", err)
				continue
			}

			if decoded.ErrorCode != internal.NO_ERROR {
				fmt.Printf("ERROR: Server responeded with Error code %d\n", decoded.ErrorCode)
				continue
			}

			fmt.Println("Success!")

		case "DELETE":
			if len(splitLine) != 2 {
				fmt.Printf("DELETE command takes one argument. Got %d.\n", len(splitLine)-1)
				continue
			}

			deleteCommand := internal.DeleteCommand{Key: splitLine[1]}
			encoded, err := deleteCommand.Encode()
			if err != nil {
				fmt.Printf("ERROR: Failed to encode DELETE command. Command: '%+v'. Error: %v\n", deleteCommand, err)
				continue
			}

			res, err := tcpConn.SendMessage(encoded)
			if err != nil {
				fmt.Printf("ERROR: Failed to send DELETE command. Command: '%+v'. Error: %v\n", deleteCommand, err)
				continue
			}

			decoded, err := internal.DecodeResponse(res)
			if err != nil {
				fmt.Printf("ERROR: Failed to decode response. Error: %v\n", err)
				continue
			}

			if decoded.ErrorCode != internal.NO_ERROR {
				fmt.Printf("ERROR: Server responeded with Error code %d\n", decoded.ErrorCode)
				continue
			}

			fmt.Println("Success!")

		case "HELP":
			fmt.Print(HELP_MESSAGE)

		default:
			fmt.Printf("Unknown Command '%s'\n%s", command, HELP_MESSAGE)
		}

		fmt.Print("> ")
	}
}
