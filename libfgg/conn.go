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
	Send([]byte, bool) error
}
