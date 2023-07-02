package httpd

import (
	"math/rand"
	"time"

	"github.com/a-wing/lightcable"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Config struct {
	Listen    string

	PathPrefix  string
	StoragePath string
}

type Server struct {
	lcSrv *lightcable.Server
	cfg   *Config
}

type MessageHello struct {
	Room string `json:"room"`
	Name string `json:"name"`
}

func NewServer(lcSrv *lightcable.Server, cfg *Config) *Server {
	return &Server{
		lcSrv: lcSrv,
		cfg:   cfg,
	}
}
