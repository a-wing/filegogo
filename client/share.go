package client

import (
	"net/url"
	"strings"
)

func IsShareInit(addr string) bool {
	if u, err := url.Parse(addr); err != nil {
		return false
	} else {
		if arr := strings.Split(u.Path, "/"); len(arr) > 2 {
			if arr[2] != "" {
				return true
			}
			return false
		}
		return false
	}
}

// http://localhost:8033/t/1024"
// To:
// ws://localhost:8033/topic/1024"
func ShareToWebSocket(addr string) string {
	if u, err := url.Parse(addr); err != nil {
		return addr
	} else {
		if u.Scheme == "https" {
			u.Scheme = "wss"
		} else {
			u.Scheme = "ws"
		}

		if arr := strings.Split(u.Path, "/"); len(arr) > 2 {
			u.Path = "/topic/" + arr[2]
		}
		return u.String()
	}
}

// ws://localhost:8033/topic/1024"
// To:
// http://localhost:8033/t/1024"
func WebSocketToShare(addr string) string {
	if u, err := url.Parse(addr); err != nil {
		return addr
	} else {
		if u.Scheme == "wss" {
			u.Scheme = "https"
		} else {
			u.Scheme = "http"
		}

		if arr := strings.Split(u.Path, "/"); len(arr) > 2 {
			u.Path = "/t/" + arr[2]
		}
		return u.String()
	}
}
