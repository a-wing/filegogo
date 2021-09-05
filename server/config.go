package server

import (
	"github.com/pion/webrtc/v3"
)

type Config struct {
	Server     string
	IcsServers *webrtc.Configuration
}
