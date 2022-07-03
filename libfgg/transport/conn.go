package transport

type Conn interface {
	Send([]byte, []byte) error
	SetOnRecv(func([]byte, []byte))
}
