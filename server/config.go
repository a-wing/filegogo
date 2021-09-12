package server

import (
	"filegogo/server/turnd"

	"github.com/pion/webrtc/v3"
)

type Config struct {
	Turn       *turnd.Config
	Server     string
	ICEServers []webrtc.ICEServer
}
