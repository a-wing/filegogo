//go:build release

package server

import (
	"embed"
	"io/fs"
	"log"
)

//go:embed build
var build embed.FS

var dist fs.FS

func init() {
	var err error
	dist, err = fs.Sub(build, "build")
	if err != nil {
		log.Fatal(err)
	}
}
