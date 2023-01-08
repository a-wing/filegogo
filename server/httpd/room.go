package httpd

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/a-wing/lightcable"
	"github.com/gorilla/mux"
	"github.com/hashicorp/golang-lru"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Config struct {
	Listen    string
	RoomAlive int
	RoomCount int

	PathPrefix  string
	StoragePath string
}

type Server struct {
	lcSrv *lightcable.Server
	cache *lru.Cache
	mutex *sync.Mutex
	cfg   *Config
}

type MessageHello struct {
	Room string `json:"room"`
	Name string `json:"name"`
}

func NewServer(lcSrv *lightcable.Server, cfg *Config) *Server {
	lcSrv.OnConnected(func(w http.ResponseWriter, r *http.Request) (string, string, bool) {
		room := mux.Vars(r)["room"]
		name := r.URL.Query().Get("name")
		log.Printf("room name: %v\n", room)
		return room, name, true
	})

	cache, err := lru.New(cfg.RoomAlive)
	if err != nil {
		panic(err)
	}
	return &Server{
		lcSrv: lcSrv,
		cache: cache,
		mutex: &sync.Mutex{},
		cfg:   cfg,
	}
}

func (s *Server) ApplyCable(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(MessageHello{
		Room: s.newRoom(),
		Name: "",
	})
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// Concurrent needs mutex lock
func (s *Server) newRoom() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	room := ""
	hasKey := func() bool {
		_, ok := s.cache.Get(room)
		return ok
	}
	for hasKey() || room == "" {
		room = strconv.Itoa(rand.Intn(s.cfg.RoomCount))
	}
	return room
}
