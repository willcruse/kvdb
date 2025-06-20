package main

import (
	"log"

	"github.com/willcruse/kvdb/v2/server/server_src"
)

func main() {
	sb := &serversrc.MapStorageBackend{}
	sb.Init()

	opLogger := &serversrc.StringDiskLogger{}
	err := opLogger.Init()
	if err != nil {
		log.Fatalf("Failed to Init StringDiskLogger. Error: %v\n", err)
	}

	tcpListener := serversrc.TCPListener{
		Address:  ":1337",
		Storage:  sb,
		OpLogger: opLogger,
	}

	tcpListener.Listen()
}
