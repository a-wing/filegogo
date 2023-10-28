//go:build !release

package server

import "embed"

//go:embed index.html
var dist embed.FS
