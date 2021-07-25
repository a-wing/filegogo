package lightcable

import (
	"context"
	"math/rand"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func makeWsProto(s string) string {
	return "ws" + strings.TrimPrefix(s, "http")
}

func TestServer(t *testing.T) {
	server := NewServer()
	ctx, cancel := context.WithCancel(context.Background())
	go server.Run(ctx)

	httpServer := httptest.NewServer(server)

	cable := "/test"
	ws, _, err := websocket.DefaultDialer.DialContext(ctx, makeWsProto(httpServer.URL+cable), nil)
	if err != nil {
		t.Error(err)
	}

	ws2, _, err := websocket.DefaultDialer.DialContext(ctx, makeWsProto(httpServer.URL+cable), nil)
	if err != nil {
		t.Error(err)
	}

	data := make([]byte, 4096)
	for i := 0; i < 10; i++ {
		n, err := rand.Read(data)
		if err != nil {
			t.Error(err)
		}
		ws.WriteMessage(websocket.TextMessage, data[:n])

		typ, recv, err := ws2.ReadMessage()
		if err != nil {
			t.Error(err)
		}

		if typ != websocket.TextMessage {
			t.Error("Type should TextMessage")
		}

		if string(recv) != string(data[:n]) {
			t.Error("Data should Equal")
		}
	}

	if err := ws.SetReadDeadline(time.Now().Add(time.Millisecond)); err != nil {
		t.Error(err)
	}

	if _, _, err := ws.ReadMessage(); err == nil {
		t.Error("Should have error")
	}

	cancel()
}
