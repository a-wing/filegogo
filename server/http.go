package server

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"filegogo/lightcable"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/hashicorp/golang-lru"
)

const (
	PrefixShare = "share"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Server struct {
	lcSrv *lightcable.Server
	cache *lru.Cache
}

type MessageHello struct {
	Share string `json:"share"`
	Token string `json:"token"`
}

func NewServer(lcSrv *lightcable.Server) *Server {
	cache, err := lru.New(1024)
	if err != nil {
		panic(err)
	}
	return &Server{
		lcSrv: lcSrv,
		cache: cache,
	}
}

func (s *Server) ApplyCable(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(MessageHello{
		Share: s.uniqueID(""),
		Token: "",
	})
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (s *Server) uniqueID(key string) string {
	if _, ok := s.cache.Get(key); ok || key == "" {
		return s.uniqueID(strconv.Itoa(rand.Intn(10000)))
	}
	return key
}

func (s *Server) JoinCable(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["id"]
	token := r.URL.Query().Get("token")

	log.Printf("topic name: %v\n", name)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	s.lcSrv.JoinCable(name, token, conn)
}
