package share

import (
	"net/url"
	"strings"

	"filegogo/server"
	"filegogo/util"
)

func IsShareInit(addr string) bool {
	if u, err := url.Parse(addr); err != nil {
		return false
	} else {
		if arr := strings.Split(u.Path, "/"); len(arr) > 1 {
			if arr[1] != "" {
				return true
			}
			return false
		}
		return false
	}
}

// http://localhost:8033/1024"
// To:
// ws://localhost:8033/<Prefix>/1024"
func ShareToWebSocket(addr string) string {
	addr = util.ProtoHttpToWs(addr)
	if u, err := url.Parse(addr); err != nil {
		return addr
	} else {
		if arr := strings.Split(u.Path, "/"); len(arr) > 1 {
			u.Path = server.Prefix + "/" + arr[1]
		}
		return u.String()
	}
}

// ws://localhost:8033/<Prefix>/1024"
// To:
// http://localhost:8033/1024"
func WebSocketToShare(addr string) string {
	addr = util.ProtoWsToHttp(addr)
	if u, err := url.Parse(addr); err != nil {
		return addr
	} else {
		if arr := strings.Split(u.Path, "/"); len(arr) > 2 {
			u.Path = arr[2]
		}
		return u.String()
	}
}
