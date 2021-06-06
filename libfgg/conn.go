package libfgg

const (
	TypeStr = true
	TypeBin = false
)

type Conn interface {
	Send([]byte, bool) error
}
