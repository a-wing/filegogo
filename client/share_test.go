package client

import "testing"

func TestIsShareInit(t *testing.T) {
	initAddr := "http://localhost:8080"
	if IsShareInit(initAddr) {
		t.Error("Not init")
	}

	joinAddr := "http://localhost:8080/t/1234"
	if !IsShareInit(joinAddr) {
		t.Error("Need init")
	}
}

func TestShareToWebSocket(t *testing.T) {
	req := "http://localhost:8033/t/1024"
	res := "ws://localhost:8033/topic/1024"

	if ShareToWebSocket(req) != res {
		t.Error("Should equal")
	}
}

func TestWebSocketToShare(t *testing.T) {
	req := "ws://localhost:8033/topic/1024"
	res := "http://localhost:8033/t/1024"

	if WebSocketToShare(req) != res {
		t.Error("Should equal")
	}
}
