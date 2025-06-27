package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/willcruse/kvdb/v2/server/server_src"
)

const (
	DEFAULT_PORT          = 1337
	DEFAULT_LOG_FILE_PATH = "kv.db"
)

type Config struct {
	Port        int
	LogFilePath string
	Help        bool
}

// Basic argument parser
// Fails on some uses e.g.
// --log-file-path --port will use '--port' as the filename
func configFromArgs(args []string) (Config, error) {
	config := Config{Port: DEFAULT_PORT, LogFilePath: DEFAULT_LOG_FILE_PATH, Help: false}

	// First arg is binary path
	for i := 1; i < len(args); i++ {
		arg := args[i]
		if arg[0] != '-' {
			// not an arg to parse. skip...
			// we could error instead
			continue
		}

		cleanArg := strings.TrimLeft(arg, "-")
		switch cleanArg {
		case "port":
			i++
			if i >= len(args) {
				return config, fmt.Errorf("(config-parsing) Expected port number to follow --port option. Did you specify a port number?")
			}
			portArg := args[i]
			portNum, err := strconv.Atoi(portArg)
			if err != nil {
				return config, fmt.Errorf("(config-parsing) Failed to parse port number from %s. Error: %+v", portArg, err)
			}
			config.Port = portNum
		case "log-file-path":
			i++
			if i >= len(args) {
				return config, fmt.Errorf("(config-parsing) Expected filepath to follow --log-file-path option. Did you add a filepath?")
			}
			logFilePath := args[i]
			config.LogFilePath = logFilePath
		case "help":
			config.Help = true
		default:
			log.Printf("Warning: Unknown arg %s\n", arg)

		}
	}

	return config, nil
}

func main() {
	config, err := configFromArgs(os.Args)
	if err != nil {
		log.Fatalf("Failed to parse CLI args. Error = %+v\n", err)
	}

	if config.Help {
		log.Println("KVDB\nOptions:\n--help: display this message and exit\n--port <INT> Port to run the server on\n--log-file-path <FILE_PATH> Filepath to store the write log to")
		os.Exit(0)
	}

	sb := &serversrc.MapStorageBackend{}
	sb.Init()

	opLogger := &serversrc.StringDiskLogger{FileName: config.LogFilePath}
	err = opLogger.Init()
	if err != nil {
		log.Fatalf("Failed to Init StringDiskLogger. Error: %v\n", err)
	}
	defer opLogger.Close()

	serverAddress := fmt.Sprintf(":%d", config.Port)
	tcpListener := serversrc.TCPListener{
		Address:  serverAddress,
		Storage:  sb,
		OpLogger: opLogger,
	}

	log.Printf("Starting server on %s\n", serverAddress)
	tcpListener.Listen()
}
