package libfgg

import (
	"context"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

type Fgg struct {
	Server string
	Action string
}

func (f *Fgg) Topic() string {
	return f.Server + "/topic/"
}

func (f *Fgg) OnShare(addr string) {
	log.Println("=== WebSocket Connected ===")
	fmt.Println(addr)
	log.Println("=== =================== ===")
}

func (f *Fgg) Send(ctx context.Context, list []string) {
	if len(list) == 0 {
		panic("Need File")
	}

	ws := NewWebSocketConn()
	ws.OnOpen = func() {
		f.OnShare(WebSocketToShare(ws.Server))
	}
	ctx, cancel := context.WithCancel(ctx)
	ws.Start(ctx, f.Topic())
	go ws.Run(ctx)

	file, err := os.Open(list[0])
	if err != nil {
		panic(err)
	}
	transfer := &Transfer{
		Conn: ws,
		File: file,
	}
	transfer.Send()
	transfer.Run()
	cancel()
}

func (f *Fgg) Recv(ctx context.Context, list []string) {
	var file *os.File
	var err error
	if len(list) != 0 {
		file, err = os.Create(list[0])
		if err != nil {
			panic(err)
		}
	}

	ws := NewWebSocketConn()
	ws.OnOpen = func() {
		f.OnShare(WebSocketToShare(ws.Server))
	}
	ctx, cancel := context.WithCancel(ctx)
	ws.Start(ctx, f.Topic())
	go ws.Run(ctx)

	transfer := &Transfer{
		Conn: ws,
		File: file,
	}
	transfer.Recv()
	transfer.Run()
	cancel()
}
