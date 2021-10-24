package cmd

import (
	"os"

	"filegogo/version"

	"github.com/urfave/cli/v2"

	log "github.com/sirupsen/logrus"
)

var (
	app = &cli.App{
		Name: "filegogo",
		Version: version.Version + " " + version.Commit + " " + version.Date,
		Usage: "A file transfer tool that can be used in the browser webrtc p2p",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "filegogo.toml",
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
