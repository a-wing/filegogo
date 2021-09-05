package client

import (
	"context"
	"fmt"

	"filegogo/client/api"
	"filegogo/client/qrcode"
	"filegogo/client/share"
	"filegogo/client/util"
	"filegogo/libfgg"
	"filegogo/libfgg/transfer"
	"filegogo/server"

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

func (t *Client) OnShare(addr string) {
	log.Println("=== Please use this address ===")

	// Show QRcode
	if t.Config.ShowQRcode {
		fmt.Println()
		qrcode.ShowQRcode(addr, t.Config.QRcodeConfig)
		fmt.Println()
	}

	fmt.Println(share.LinksToShare(addr))
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

func (c *Client) overrideServer() {
	if !share.IsShareInit(c.Config.Server) {
		room, err := api.GetRoom(c.Config.Server + server.Prefix)
		if err != nil {
			log.Debug(c.Config.Server + server.Prefix)
			panic(err)
		}
		c.Config.Server += server.Prefix + room
		c.OnShare(c.Config.Server)
	} else {
		c.Config.Server = share.ShareToLinks(c.Config.Server)
	}
}

func (c *Client) Send(ctx context.Context, files []string) {
	fgg := libfgg.NewFgg()
	fgg.Tran.OnProgress = c.OnProgress
	fgg.OnPreTran = c.OnPreTran

	c.overrideServer()

	fgg.UseWebsocket(util.ProtoHttpToWs(c.Config.Server))
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

	c.overrideServer()

	fgg.UseWebsocket(util.ProtoHttpToWs(c.Config.Server))
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
