package cmd

import (
	"os"

	"github.com/urfave/cli/v2"

	log "github.com/sirupsen/logrus"
)

var (
	app = &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Value: false,
				Usage: "Enabled Debug mode",
			},
		},
		Before: func(c *cli.Context) error {
			if c.Bool("debug") {
				log.SetReportCaller(true)
				log.SetLevel(log.DebugLevel)
			}
			log.Warnln(c.Path("config"))
			log.Debugln(c.Path("config"))
			return nil
		},
	}
)

func Execute() {
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
