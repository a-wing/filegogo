package cmd

import (
	"context"

	"filegogo/client"

	"github.com/urfave/cli/v2"
)

func init() {
	app.Commands = append(app.Commands, &cli.Command{
		Name:  "recv",
		Usage: "recv <file>",
		Action: func(c *cli.Context) error {
			config := client.DefaultConfig
			config.Server = c.String("share")
			cli, err := client.NewClient(config)
			if err != nil {
				panic(err)
			}

			cli.Recv(context.Background(), c.Args().Slice())
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
