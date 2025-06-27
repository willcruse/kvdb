package listener

import "net"

// Listener is responsible for responding to incoming messages
type Listener interface {
	Listen(connChan chan net.Conn) error
}
