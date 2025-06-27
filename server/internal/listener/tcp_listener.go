package listener

import (
	"log"
	"net"
)

type TCPListener struct {
	Address string
}

func (t *TCPListener) Listen(connChan chan Readable) error {
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

		connChan <- conn

		// go t.handleNetConn(conn)

	}
}
