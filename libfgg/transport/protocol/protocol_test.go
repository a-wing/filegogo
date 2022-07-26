package protocol

import (
	"testing"
)

func TestProtocol(t *testing.T) {
	head, body := []byte("hello"), []byte("world world !!")

	data := Encode(head, body)

	head2, body2 := Decode(data)
	if string(head2) != string(head) {
		t.Error(head)
	}
	if string(body2) != string(body) {
		t.Error(body)
	}
}
