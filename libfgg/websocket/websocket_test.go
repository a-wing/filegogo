package websocket

import (
	"testing"
)

func TestFillAddr(t *testing.T) {
	result := "ws://localhost:8080/share/"
	if r := fillAddr("ws://localhost:8080"); r != result {
		t.Error(r)
	}

	if r := fillAddr("ws://localhost:8080/"); r != result {
		t.Error(r)
	}

	if r := fillAddr("ws://localhost:8080/xxx"); r != result {
		t.Error(r)
	}

	if r := fillAddr("ws://localhost:8080/share"); r != result {
		t.Error(r)
	}

	if r := fillAddr("ws://localhost:8080/share/1234"); r != result + "1234" {
		t.Error(r)
	}
}

func TestServer(t *testing.T) {
	c := NewConn("ws://localhost:8080/")
	c.updateServer("1234")

	if addr := c.authServer(); addr != "ws://localhost:8080/share/1234" {
		t.Error(addr)
	}

	c.token = "xxx"
	if addr := c.authServer(); addr != "ws://localhost:8080/share/1234?token=xxx" {
		t.Error(addr)
	}
}

