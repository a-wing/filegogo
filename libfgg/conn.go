package libfgg

const (
	TextMessage   = 1
	BinaryMessage = 2
)

type Message struct {
	Type int
	Data []byte
}

type Conn interface {
	Send(t int, data []byte) error
	Recv() (t int, data []byte, err error)
}
