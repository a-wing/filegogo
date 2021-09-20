package client

import (
	"context"
	"fmt"

	"filegogo/client/api"
	"filegogo/client/qrcode"
	"filegogo/client/util"
	"filegogo/libfgg"
	"filegogo/libfgg/transfer"
	"filegogo/server"

	bar "github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
)

type ClientConfig struct {
	Server string

	ShowQRcode   bool
	ShowProgress bool
	ServerConfig *server.ApiConfig
	QRcodeConfig *qrcode.Config
	Level        string
}

type Client struct {
	Config *ClientConfig
	api    *api.Api
	bar    *bar.ProgressBar
}

func NewClient(config *ClientConfig) (*Client, error) {
	return &Client{
		api:    api.NewApi(config.Server),
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

func (c *Client) overrideServer() {
	if cfg, err := c.api.GetConfig(); err != nil {
		panic(err)
	} else {
		// local iceServers merge remote iceServers
		c.Config.ServerConfig.ICEServers = append(c.Config.ServerConfig.ICEServers, cfg.ICEServers...)
	}

	if !c.api.HasRoom() {
		_, err := c.api.NewRoom()
		if err != nil {
			log.Debug(c.api.RoomAddress())
			panic(err)
		}
		c.OnShare(c.api.ToShare())
	}
}

func (c *Client) Send(ctx context.Context, files []string) {
	fgg := libfgg.NewFgg()
	fgg.Tran.OnProgress = c.OnProgress
	fgg.OnPreTran = c.OnPreTran

	c.overrideServer()

	fgg.UseWebsocket(util.ProtoHttpToWs(c.api.RoomAddress()))
	if err := fgg.Send(files); err != nil {
		panic(err)
	}
	fgg.UseWebRTC(c.Config.ServerConfig.ICEServers)
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

	fgg.UseWebsocket(util.ProtoHttpToWs(c.api.RoomAddress()))
	if err := fgg.Recv(files); err != nil {
		panic(err)
	}
	fgg.UseWebRTC(c.Config.ServerConfig.ICEServers)
	if err := fgg.Run(); err != nil {
		fmt.Println()
		fmt.Println(err)
	} else {
		fmt.Println()
		fmt.Println("md5 VerifyHash successful")
	}
}
