package transport

type Transport struct {
	conns []Conn
}

func New() *Transport {
	return &Transport{}
}
