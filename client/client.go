package client

import (
	"context"
	"fmt"

	"filegogo/client/api"
	"filegogo/client/qrcode"
	"filegogo/client/util"
	"filegogo/libfgg"
	"filegogo/libfgg/pool"
	"filegogo/server/config"

	bar "github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
)

var (
	DefaultConfig = &ClientConfig{
		ShowQRcode:   true,
		ShowProgress: true,
		ServerConfig: &config.ApiConfig{},
		QRcodeConfig: &qrcode.Config{
			Foreground: "black",
			Background: "white",
			Level:      "low",
			Align:      "left",
		},
		NoIceServer: false,
	}
)

type ClientConfig struct {
	Server string

	ShowQRcode   bool
	ShowProgress bool
	ServerConfig *config.ApiConfig
	QRcodeConfig *qrcode.Config
	Level        string
	NoIceServer  bool
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

func (t *Client) OnPreTran(file *pool.Meta) {
	if t.Config.ShowProgress {
		t.bar = bar.New64(int64(file.Size))
	}
}

func (t *Client) OnProgress(c int64) {
	if t.Config.ShowProgress {
		t.bar.Set64(c)
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
	fgg.OnPreTran = c.OnPreTran
	fgg.SetOnProgress(c.OnProgress)

	// TODO:
	//ctx, cancel := context.WithCancel(ctx)
	//fgg.OnPostTran = func(h *pool.Hash) {
	//	cancel()
	//}

	c.overrideServer()

	fgg.UseWebsocket(util.ProtoHttpToWs(c.api.RoomAddress()))
	if err := fgg.SetSend(files[0]); err != nil {
		panic(err)
	}

	if !c.Config.NoIceServer {
		log.Println("use webrtc")
		fgg.UseWebRTC(c.Config.ServerConfig.ICEServers)
	}

	select {
	case <-ctx.Done():
	}
}

func (c *Client) Recv(ctx context.Context, files []string) {
	fgg := libfgg.NewFgg()
	fgg.OnPreTran = c.OnPreTran
	fgg.SetOnProgress(c.OnProgress)

	ctx, cancel := context.WithCancel(ctx)
	fgg.OnPostTran = func(h *pool.Hash) {
		cancel()
	}
	c.overrideServer()

	fgg.UseWebsocket(util.ProtoHttpToWs(c.api.RoomAddress()))
	if err := fgg.SetRecv(files[0]); err != nil {
		panic(err)
	}
	if !c.Config.NoIceServer {
		log.Println("use webrtc")
		fgg.UseWebRTC(c.Config.ServerConfig.ICEServers)
	}

	ch := make(chan bool)
	fgg.OnRecvFile = func(meta *pool.Meta) {
		ch <- true
	}

	go fgg.GetMeta()

	<-ch
	log.Println("start download file")
	if err := fgg.Run(ctx); err != nil {
		log.Println()
		log.Println(err)
	} else {
		log.Println()
		log.Println("md5 VerifyHash successful")
	}
}
