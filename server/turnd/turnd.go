package turnd

import (
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"

	"github.com/hashicorp/golang-lru"
	"github.com/pion/turn/v2"
)

type Config struct {
	User         string
	Realm        string
	Listen       string
	PublicIP     string
	RelayMinPort int
	RelayMaxPort int
}

type Server struct {
	cfg      *Config
	usersMap *lru.Cache
}

func New(cfg *Config) *Server {
	usersMap, err := lru.New(1024)
	if err != nil {
		panic(err)
	}
	return &Server{
		cfg:      cfg,
		usersMap: usersMap,
	}
}

func RandomUser() (string, string) {
	return strconv.Itoa(rand.Intn(1000000)), strconv.Itoa(rand.Intn(1000000))
}

func (s *Server) NewUser(user string) {
	username := strings.Split(user, ":")[0]
	password := strings.Split(user, ":")[1]
	log.Println("Add Turn Server User:", username, password)
	s.usersMap.Add(username, turn.GenerateAuthKey(username, s.cfg.Realm, password))
}

func (s *Server) Run() (*turn.Server, error) {
	udpListener, err := net.ListenPacket("udp4", s.cfg.Listen)
	if err != nil {
		return nil, err
	}

	if user := s.cfg.User; user != "" {
		s.NewUser(user)
	}

	return turn.NewServer(turn.ServerConfig{
		Realm: s.cfg.Realm,
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
					Address:      "0.0.0.0",                   // But actually be listening on every interface
					MinPort:      uint16(s.cfg.RelayMinPort),
					MaxPort:      uint16(s.cfg.RelayMaxPort),
				},
			},
		},
	})
}
