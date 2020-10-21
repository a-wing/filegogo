// +build dev

package data

import (
	"net/http"
)

var (
	Dir = http.Dir("dist/")
)
