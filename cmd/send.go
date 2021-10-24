package cmd

import (
	"context"

	"filegogo/client"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli/v2"
)

func init() {
	config := client.DefaultConfig
	app.Commands = append(app.Commands, &cli.Command{
		Name:  "send",
		Usage: "send <file>",
		Before: func(c *cli.Context) error {
			toml.DecodeFile(c.Path("config"), config)
			config.Server = c.String("share")
			return nil
		},
		Action: func(c *cli.Context) error {
			cli, err := client.NewClient(config)
			if err != nil {
				panic(err)
			}

			cli.Send(context.Background(), c.Args().Slice())
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "share",
				Aliases: []string{"s"},
				Value:   "https://send.22333.fun",
				Usage:   "Share Link",
			},
		},
	})
}
