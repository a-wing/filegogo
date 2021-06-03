package client

import (
	"context"
	"fmt"
	"os"

	fgg "filegogo/libfgg"
	"filegogo/client/qrcode"

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
	qrcode.ShowQRcode(addr, nil)
	fmt.Println(addr)
	log.Println("=== =================== ===")
}

func (t *Client) OnPreTran(file *fgg.MetaFile) {
	t.bar = bar.New64(file.Size)
}

func (f *Client) Send(ctx context.Context, list []string) {
	if len(list) == 0 {
		panic("Need File")
	}

	ws := fgg.NewWebSocketConn()
	ws.OnOpen = func() {
		f.OnShare(WebSocketToShare(ws.Server))
	}
	ctx, cancel := context.WithCancel(ctx)
	ws.Start(ctx, ShareToWebSocket(f.Topic()))
	go ws.Run(ctx)

	file, err := os.Open(list[0])
	if err != nil {
		panic(err)
	}
	transfer := &fgg.Fgg{
		Conn: ws,
		File: file,
		OnPreTran: func(fl *fgg.MetaFile) {
			f.OnPreTran(fl)
		},
	}
	transfer.Send()
	transfer.Tran.OnProgress = func(c int64) {
		f.bar.Add64(c)
	}
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
		f.OnShare(WebSocketToShare(ws.Server))
	}
	ctx, cancel := context.WithCancel(ctx)
	ws.Start(ctx, ShareToWebSocket(f.Topic()))
	go ws.Run(ctx)

	transfer := &fgg.Fgg{
		Conn: ws,
		File: file,
		OnPreTran: func(fl *fgg.MetaFile) {
			f.OnPreTran(fl)
		},
	}
	transfer.Recv()
	transfer.Tran.OnProgress = func(c int64) {
		f.bar.Add64(c)
	}
	transfer.Run()
	cancel()
}
