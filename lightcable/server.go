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

func init() {
	rand.Seed(time.Now().UnixNano())
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func CreateTopic(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	name := strconv.Itoa(rand.Intn(10000))
	token := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", &conn)))
	hub.Cache.Add(token, &conn)
	topic := NewTopic(name, conn)
	hub.Add(name, topic)

	if err := conn.WriteJSON(&struct {
		Topic string `json:"topic"`
		Token string `json:"token"`
	}{
		Topic: name,
		Token: token,
	}); err != nil {
		log.Println(err)
	}

	go readPump(hub, topic, conn)
}

func JoinTopic(hub *Hub, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("topic ID: %v\n", vars["id"])
	id := vars["id"]

	if topic := hub.Topic[id]; topic != nil {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
		}
		topic.Register(conn)
		go readPump(hub, topic, conn)
	} else {
		if _, ok := hub.Cache.Get(r.URL.Query().Get("token")); ok {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Println(err)
			}
			topic := NewTopic(id, conn)
			hub.Add(id, topic)
			go readPump(hub, topic, conn)
		}
	}
}
