package client

import (
	"context"
	"fmt"

	"filegogo/client/api"
	"filegogo/client/qrcode"
	"filegogo/libfgg"
	"filegogo/libfgg/transfer"
	"filegogo/util"

	"github.com/pion/webrtc/v3"

	bar "github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
)

type ClientConfig struct {
	Server string

	ShowQRcode   bool
	ShowProgress bool
	IcsServers   *webrtc.Configuration
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
	log.Println("=== Please use this address ===")

	// Show QRcode
	if t.Config.ShowQRcode {
		fmt.Println()
		qrcode.ShowQRcode(addr, t.Config.QRcodeConfig)
		fmt.Println()
	}

	fmt.Println(addr)
	log.Println("=== ======================= ===")
}

func (t *Client) OnPreTran(file *transfer.MetaFile) {
	if t.Config.ShowProgress {
		t.bar = bar.New64(file.Size)
	}
}

func (t *Client) OnProgress(c int64) {
	if t.Config.ShowProgress {
		t.bar.Add64(c)
	}
}

func (c *Client) Send(ctx context.Context, files []string) {
	fgg := libfgg.NewFgg()
	fgg.Tran.OnProgress = c.OnProgress
	fgg.OnPreTran = c.OnPreTran

	room, err := api.GetRoom(c.Config.Server)
	if err != nil {
		panic(err)
	}
	c.OnShare(c.Config.Server+room)

	fgg.UseWebsocket(util.ProtoHttpToWs(c.Config.Server)+room)
	if err := fgg.Send(files); err != nil {
		panic(err)
	}
	fgg.UseWebRTC(c.Config.IcsServers)
	if err := fgg.Run(); err != nil {
		fmt.Println()
		fmt.Println(err)
	} else {
		fmt.Println()
	}
}

func (c *Client) Recv(ctx context.Context, files []string) {
	fgg := libfgg.NewFgg()
	fgg.Tran.OnProgress = c.OnProgress
	fgg.OnPreTran = func(t *transfer.MetaFile) {
		c.OnPreTran(t)
		go func() {
			fgg.RunWebRTC()
			fgg.GetFile()
		}()
	}

	room, err := api.GetRoom(c.Config.Server)
	if err != nil {
		panic(err)
	}

	fgg.UseWebsocket(util.ProtoHttpToWs(c.Config.Server)+room)
	if err := fgg.Recv(files); err != nil {
		panic(err)
	}
	fgg.UseWebRTC(c.Config.IcsServers)
	if err := fgg.Run(); err != nil {
		fmt.Println()
		fmt.Println(err)
	} else {
		fmt.Println()
		fmt.Println("md5 VerifyHash successful")
	}
}
