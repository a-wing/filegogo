package turnd

import (
	"net"

	"github.com/hashicorp/golang-lru"
	"github.com/pion/turn/v2"
)

type Config struct {
	Username string
	Password string
	Listen string
	Realm string
	PublicIP string
	RelayMinPort int
	RelayMaxPort int
}

type Server struct {
	cfg *Config
	usersMap *lru.Cache
}

func New(cfg *Config) *Server {
	usersMap, err := lru.New(1024)
	if err != nil {
		panic(err)
	}
	return &Server{
		cfg: cfg,
		usersMap: usersMap,
	}
}

func (s *Server) NewUser(user string) {
	s.usersMap.Add(user, turn.GenerateAuthKey(s.cfg.Username, s.cfg.Realm, s.cfg.Password))
}

func (s *Server) Run() (*turn.Server, error) {
	udpListener, err := net.ListenPacket("udp4", s.cfg.Listen)
	if err != nil {
		return nil, err
	}

	return turn.NewServer(turn.ServerConfig{
		Realm:         s.cfg.Realm,
		AuthHandler: func(username, realm string, srcAddr net.Addr) (key []byte, ok bool) {
			if key, ok := s.usersMap.Get(username); ok {
				return key.([]byte), true
			}
			return nil, false
		},
		PacketConnConfigs: []turn.PacketConnConfig{
			{
				PacketConn: udpListener,
				RelayAddressGenerator: &turn.RelayAddressGeneratorPortRange{
					RelayAddress: net.ParseIP(s.cfg.PublicIP), // Claim that we are listening on IP passed by user (This should be your Public IP)
					Address:      "0.0.0.0",                 // But actually be listening on every interface
					MinPort:      uint16(s.cfg.RelayMinPort),
					MaxPort:      uint16(s.cfg.RelayMaxPort),
				},
			},
		},
	})
}
