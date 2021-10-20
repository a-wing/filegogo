package cmd

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
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
