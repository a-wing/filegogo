package socket

import (
	"context"
	"net"
	"testing"
	"time"
)

func TestSocket(t *testing.T) {
	client, server := net.Pipe()

	head, body := []byte("hello"), []byte("world world !!")

	cc := New(client)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	go cc.Run(ctx)

	sign := make(chan bool)
	cs := New(server)
	cs.SetOnRecv(func(h, b []byte) {
		if string(h) != string(head) {
			t.Error(h)
		}
		if string(b) != string(body) {
			t.Error(b)
		}
		sign <- true
	})
	go cs.Run(ctx)

	cc.Send(head, body)

	<-sign
}
