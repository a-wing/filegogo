package cmd

import (
	"filegogo/server"
	"filegogo/server/httpd"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli/v2"

	log "github.com/sirupsen/logrus"
)

func init() {
	config := &server.Config{
		Http: &httpd.Config{
			Listen:    "0.0.0.0:8080",
			RoomAlive: 1024,
			RoomCount: 10000,
		},
	}
	app.Commands = append(app.Commands, &cli.Command{
		Name:  "server",
		Usage: "websocket broker server",
		Before: func(c *cli.Context) error {
			toml.DecodeFile(c.Path("config"), config)
			log.Debugln(config)
			return nil
		},
		Action: func(c *cli.Context) error {
			server.Run(config)
			return nil
		},
	})
}
