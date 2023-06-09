package config

import (
	"filegogo/server/httpd"
	"filegogo/server/turnd"

	"github.com/pion/webrtc/v3"
)

type Config struct {
	Http       *httpd.Config
	Turn       *turnd.Config
	ICEServers []webrtc.ICEServer
}

type ApiConfig struct {
	ICEServers []webrtc.ICEServer `json:"iceServers,omitempty"`
}
