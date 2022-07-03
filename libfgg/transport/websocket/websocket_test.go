package websocket

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/a-wing/lightcable"
	"github.com/gorilla/websocket"
)

func makeWsProto(s string) string {
	return "ws" + strings.TrimPrefix(s, "http")
}

func makeConns(t testing.TB, server http.Handler, rooms ...string) []*websocket.Conn {
	httpServer := httptest.NewServer(server)
	conns := make([]*websocket.Conn, len(rooms))
	var err error
	for i, room := range rooms {
		if conns[i], _, err = websocket.DefaultDialer.Dial(makeWsProto(httpServer.URL+room), nil); err != nil {
			t.Error(err)
		}
	}
	return conns
}

func TestWebsocket(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	broker := lightcable.New(lightcable.DefaultConfig)
	go broker.Run(ctx)
	conns := makeConns(t, broker, "/test", "/test")
	client, server := conns[0], conns[1]

	head, body := []byte("hello"), []byte("world world !!")

	cc := New(client)

	defer cancel()

	time.Sleep(time.Second)

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

	if err := cc.Send(head, body); err != nil {
		t.Error(err)
	}

	<-sign
}
