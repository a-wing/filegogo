package client

import (
	"context"
	"fmt"
	"os"

	"filegogo/client/qrcode"
	"filegogo/libfgg"

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
	bar    *bar.ProgressBar
}

func NewClient(config *ClientConfig) (*Client, error) {
	return &Client{
		Config: config,
	}, nil
}

func (c *Client) Topic() string {
	return c.Config.Server + "/topic/"
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

func (t *Client) OnPreTran(file *libfgg.MetaFile) {
	if t.Config.Progress {
		t.bar = bar.New64(file.Size)
	}
}

func (t *Client) OnProgress(c int64) {
	if t.Config.Progress {
		t.bar.Add64(c)
	}
}

func (c *Client) Send(ctx context.Context, list []string) {
	if len(list) == 0 {
		panic("Need File")
	}

	file, err := os.Open(list[0])
	if err != nil {
		panic(err)
	}

	fgg := libfgg.NewFgg(file)
	fgg.OnShare = c.OnShare
	fgg.Tran.OnProgress = c.OnProgress
	fgg.OnPreTran = func(fl *libfgg.MetaFile) {
		c.OnPreTran(fl)
	}

	fgg.Start(ShareToWebSocket(c.Config.Server))
	fgg.Send()
	fgg.Run()
}

func (c *Client) Recv(ctx context.Context, list []string) {
	var file *os.File
	var err error
	if len(list) != 0 {
		file, err = os.Create(list[0])
		if err != nil {
			panic(err)
		}
	}

	fgg := libfgg.NewFgg(file)
	fgg.OnShare = c.OnShare
	fgg.Tran.OnProgress = c.OnProgress
	fgg.OnPreTran = func(fl *libfgg.MetaFile) {
		c.OnPreTran(fl)
	}

	fgg.Start(ShareToWebSocket(c.Config.Server))
	fgg.Recv()
	fgg.Run()
}
