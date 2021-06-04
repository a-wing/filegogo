package lightcable

import (
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const (
	PrefixShare = "share"
	PrefixShort = "s"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func JoinTopic(hub *Hub, w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["id"]
	topic := hub.Topic[name]
	token := r.URL.Query().Get("token")
	if _, ok := hub.Cache.Get(token); !ok && topic == nil && name != "" {
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
		hub.Cache.Add(token, &conn)
	}

	// Topic Register websocket.conn
	if topic == nil {
		topic = NewTopic(name, conn)
		hub.Add(name, topic)
	} else {
		topic.Register(conn)
	}

	// websocket response
	if err := conn.WriteJSON(&MessageHello{
		Share: name,
		Token: token,
	}); err != nil {
		log.Println(err)
	}

	go readPump(hub, topic, conn)
}

type MessageHello struct {
	Share string `json:"share"`
	Token string `json:"token"`
}
