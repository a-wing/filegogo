package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"filegogo/client"
	"filegogo/client/qrcode"
	"filegogo/server"

	"github.com/pion/webrtc/v3"
	"github.com/urfave/cli/v2"
)

const serverAddr = "http://localhost:8033/s/"

func Execute() {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.UseShortOptionHandling = true
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Load configuration from `FILE`",
		},
		&cli.BoolFlag{
			Name:    "qrcode",
			Aliases: []string{"q"},
			Usage:   "Show QRcode",
			Value:   false,
		},
		&cli.BoolFlag{
			Name:    "progress",
			Aliases: []string{"p"},
			Usage:   "Show Progress Bar",
			Value:   false,
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:  "server",
			Usage: "websocket broker server",
			Action: func(c *cli.Context) error {
				fmt.Println("listen: ", c.String("listen"))
				server.Run(c.String("listen"), "TODO")
				return nil
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "listen",
					Aliases: []string{"l"},
					Value:   "0.0.0.0:8033",
					Usage:   "set server listen address and port",
				},
			},
		},
		{
			Name:  "send",
			Usage: "Send File",
			Action: func(c *cli.Context) error {
				fmt.Println("send: ", c.Args().First())
				fmt.Println("server: ", c.String("server"))
				qrconfig := &qrcode.Config{}
				iceservers := &webrtc.Configuration{}

				config := &client.ClientConfig{
					Server:       c.String("server"),
					ShowQRcode:   c.Bool("qrcode"),
					ShowProgress: c.Bool("progress"),
					IcsServers:   iceservers,
					QRcodeConfig: qrconfig,
				}
				cc, _ := client.NewClient(config)
				cc.Send(context.Background(), c.Args().Slice())
				return nil
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "server",
					Aliases: []string{"s"},
					Value:   serverAddr,
					Usage:   "Signal Server Address",
				},
			},
		},
		{
			Name:  "recv",
			Usage: "Recv File",
			Action: func(c *cli.Context) error {
				fmt.Println("recv: ", c.Args().First())
				fmt.Println("server: ", c.String("server"))
				qrconfig := &qrcode.Config{}
				iceservers := &webrtc.Configuration{}

				config := &client.ClientConfig{
					Server:       c.String("server"),
					ShowQRcode:   c.Bool("qrcode"),
					ShowProgress: c.Bool("progress"),
					IcsServers:   iceservers,
					QRcodeConfig: qrconfig,
				}
				cc, _ := client.NewClient(config)
				cc.Recv(context.Background(), c.Args().Slice())
				return nil
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "server",
					Aliases: []string{"s"},
					Value:   serverAddr,
					Usage:   "Signal Server Address",
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
