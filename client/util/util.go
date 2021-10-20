package util

import (
	"strings"
)

func ProtoHttpToWs(s string) string {
	return "ws" + strings.TrimPrefix(s, "http")
}

func ProtoWsToHttp(s string) string {
	return "http" + strings.TrimPrefix(s, "ws")
}
