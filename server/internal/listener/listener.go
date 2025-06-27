package listener

// import "net"

type Readable interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	Close() error
}

// Listener is responsible for responding to incoming messages
type Listener interface {
	Listen(connChan chan Readable) error
}
