package cmd

import (
	"filegogo/server"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func init() {
	config := &server.Config{}
	app.Commands = append(app.Commands, &cli.Command{
		Name:  "server",
		Usage: "websocket broker server",
		Before: func(c *cli.Context) error {
			toml.DecodeFile("filegogo-server.toml", config)
			return nil
		},
		Action: func(c *cli.Context) error {
			//log.Println("listen: ", c.String("listen"))
			log.Println("server")
			log.Printf("%+v\n", config)
			server.Run(config)
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
	})
}
