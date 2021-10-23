package cmd

import (
	"os"

	"github.com/urfave/cli/v2"

	log "github.com/sirupsen/logrus"
)

var (
	app = cli.NewApp()
)

func Execute() {
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
