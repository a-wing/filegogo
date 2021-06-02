package client

import (
	"context"
	"fmt"
	"os"

	fgg "filegogo/libfgg"

	bar "github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	Server string
	bar    *bar.ProgressBar
}

func (c *Client) Topic() string {
	return c.Server + "/topic/"
}

func (f *Client) OnShare(addr string) {
	log.Println("=== WebSocket Connected ===")
	fmt.Println(addr)
	log.Println("=== =================== ===")
}

func (t *Client) OnPreTran(file *fgg.FileList) {
	t.bar = bar.New64(file.Size)
}

func (f *Client) Send(ctx context.Context, list []string) {
	if len(list) == 0 {
		panic("Need File")
	}

	ws := fgg.NewWebSocketConn()
	ws.OnOpen = func() {
		f.OnShare(fgg.WebSocketToShare(ws.Server))
	}
	ctx, cancel := context.WithCancel(ctx)
	ws.Start(ctx, f.Topic())
	go ws.Run(ctx)

	file, err := os.Open(list[0])
	if err != nil {
		panic(err)
	}
	transfer := &fgg.Transfer{
		Conn: ws,
		File: file,
		OnProgress: func(c int64) {
			f.bar.Add64(c)
		},
		OnPreTran: func(fl *fgg.FileList) {
			f.OnPreTran(fl)
		},
	}
	transfer.Send()
	transfer.Run()
	cancel()
}

func (f *Client) Recv(ctx context.Context, list []string) {
	var file *os.File
	var err error
	if len(list) != 0 {
		file, err = os.Create(list[0])
		if err != nil {
			panic(err)
		}
	}

	ws := fgg.NewWebSocketConn()
	ws.OnOpen = func() {
		f.OnShare(fgg.WebSocketToShare(ws.Server))
	}
	ctx, cancel := context.WithCancel(ctx)
	ws.Start(ctx, f.Topic())
	go ws.Run(ctx)

	transfer := &fgg.Transfer{
		Conn: ws,
		File: file,
		OnProgress: func(c int64) {
			f.bar.Add64(c)
		},
		OnPreTran: func(fl *fgg.FileList) {
			f.OnPreTran(fl)
		},
	}
	transfer.Recv()
	transfer.Run()
	cancel()
}
