package protocol

import (
	"encoding/binary"
)

const (
	LengthHead = 2
	LengthBody = 2
)

func Decode(data []byte) (head, body []byte) {
	l1, l2 := binary.BigEndian.Uint16(data[:LengthHead]), binary.BigEndian.Uint16(data[LengthHead:LengthHead+LengthBody])

	payload := data[LengthHead+LengthBody:]
	return payload[:l1], payload[l1 : l1+l2]
}

func Encode(head, body []byte) []byte {
	l1, l2 := len(head), len(body)
	l1b, l2b := make([]byte, LengthHead), make([]byte, LengthBody)
	binary.BigEndian.PutUint16(l1b, uint16(l1))
	binary.BigEndian.PutUint16(l2b, uint16(l2))
	return append(append(append(l1b, l2b...), head...), body...)
}
