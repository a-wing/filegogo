package client

import (
	"context"
	"fmt"
	"os"

	"filegogo/client/qrcode"
	fgg "filegogo/libfgg"

	bar "github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
)

type ClientConfig struct {
	Server   string
	QRcode   bool
	Progress bool

	QRcodeConfig *qrcode.Config
	Level        string
}

type Client struct {
	Config *ClientConfig
	Server string
	bar    *bar.ProgressBar
}

func NewClient(config *ClientConfig) (*Client, error) {
	return &Client{
		Config: config,
	}, nil
}

func (c *Client) Topic() string {
	return c.Server + "/topic/"
}

func (t *Client) OnShare(addr string) {
	log.Println("=== WebSocket Connected ===")

	// Show QRcode
	if t.Config.QRcode {
		fmt.Println()
		qrcode.ShowQRcode(addr, t.Config.QRcodeConfig)
		fmt.Println()
	}

	fmt.Println(addr)
	log.Println("=== =================== ===")
}

func (t *Client) OnPreTran(file *fgg.MetaFile) {
	if t.Config.Progress {
		t.bar = bar.New64(file.Size)
	}
}

func (t *Client) OnProgress(c int64) {
	if t.Config.Progress {
		t.bar.Add64(c)
	}
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
	transfer.Tran.OnProgress = f.OnProgress
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
	transfer.Tran.OnProgress = f.OnProgress
	transfer.Run()
	cancel()
}
