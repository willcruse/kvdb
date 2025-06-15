package main

import "github.com/willcruse/kvdb/v2/server/server_src"

func main() {
	sb := &serversrc.MapStorageBackend{}
	sb.Init()
	tcpListener := serversrc.TCPListener{
		Address: ":1337",
		Storage: sb,
	}

	tcpListener.Listen()
}
