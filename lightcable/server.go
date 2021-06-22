package lightcable

import (
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/hashicorp/golang-lru"
)

const (
	PrefixShare = "share"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Server struct {
	topic map[string]*topic
	cache *lru.Cache
	mutex sync.Mutex

	// Register requests from the clients.
	register chan *Client

	// Inbound messages from the clients.
	broadcast chan []byte

	// Unregister requests from clients.
	unregister chan *Client
}

func NewServer() *Server {
	cache, err := lru.New(1024)
	if err != nil {
		panic(err)
	}
	return &Server{
		topic: make(map[string]*topic),
		cache: cache,

		register:   make(chan *Client),
		broadcast:  make(chan []byte),
		unregister: make(chan *Client),
	}
}

func (s *Server) JoinTopic(w http.ResponseWriter, r *http.Request) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	name := mux.Vars(r)["id"]

	//topic := hub.Topic[name]
	topic, ok := s.topic[name]
	if ok {
	}

	token := r.URL.Query().Get("token")
	if _, ok := s.cache.Get(token); !ok && topic == nil && name != "" {
		log.Printf("reject topic name: %v\n", name)
		return
	}

	// === allow websocket connect ===
	// Default topic name
	if name == "" {
		name = strconv.Itoa(rand.Intn(10000))
	}

	log.Printf("topic name: %v\n", name)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	// Default token
	if token == "" {
		token = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", &conn)))
		s.cache.Add(token, &conn)
	}

	// Topic Register websocket.conn
	if topic == nil {
		topic = NewTopic(name, s)
		go topic.run()
		s.cache.Add(name, topic)
		s.topic[name] = topic
	}
	topic.register <- &Client{topic: topic, conn: conn, send: make(chan message, 256)}

	// websocket response
	if err := conn.WriteJSON(&MessageHello{
		Share: name,
		Token: token,
	}); err != nil {
		log.Println(err)
	}
}

type MessageHello struct {
	Share string `json:"share"`
	Token string `json:"token"`
}
