package server

import (
	"filegogo/server/httpd"
	"filegogo/server/turnd"

	"github.com/pion/webrtc/v3"
)

type Config struct {
	Http       *httpd.Config
	Turn       *turnd.Config
	Server     string
	ICEServers []webrtc.ICEServer
}
