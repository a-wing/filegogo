package lightcable

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var sequence uint64
var sequenceMutex sync.Mutex

func getSequence() string {
	sequenceMutex.Lock()
	id := strconv.FormatUint(sequence, 10)
	sequence++
	sequenceMutex.Unlock()
	return id
}

func CreateTopic(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	name := getSequence()
	topic := NewTopic(name, conn)
	hub.Add(name, topic)

	msg, err := json.Marshal(&struct {
		Topic string `json:"topic"`
	}{
		Topic: name,
	})
	if err != nil {
		log.Println(err)
	}

	topic.Broadcast(&Message{
		code: websocket.TextMessage,
		data: msg,
		conn: nil,
	})

	go readPump(hub, topic, conn)
}

func JoinTopic(hub *Hub, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("topic ID: %v\n", vars["id"])
	id := vars["id"]

	if topic := hub.Topics[id]; topic != nil {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
		}
		topic.Register(conn)
		go readPump(hub, topic, conn)
	}
}

func readPump(hub *Hub, topic *Topic, conn *websocket.Conn) {
	log.Printf("Topic: %p, conn: %p opened", topic, conn)
	defer func() {
		topic.Unregister(conn)
		conn.Close()
		log.Printf("Topic: %p, conn: %p closed", topic, conn)
		if len(topic.conns) == 0 {
			hub.Remove(topic.Name)
			log.Printf("Hub: %p, closed", topic)
		}
	}()
	for {
		code, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		topic.Broadcast(&Message{
			code: code,
			data: message,
			conn: conn,
		})
	}
}
