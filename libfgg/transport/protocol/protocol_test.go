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

func TestProtocolDecodeNil(t *testing.T) {
	var data []byte

	head, body := Decode(data)
	if len(head) != 0 {
		t.Error(head)
	}
	if len(body) != 0 {
		t.Error(body)
	}
}

func TestProtocolDecodeError(t *testing.T) {
	data := []byte("world world !!")

	head, body := Decode(data)
	if len(head) != 0 {
		t.Error(head)
	}
	if len(body) != 0 {
		t.Error(body)
	}
}
