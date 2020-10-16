package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
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

func NewTopic(conn *websocket.Conn) *Cable {
	return &Cable{
		Message: &Message{
			msg:  []byte{},
			conn: conn,
		},
		conns: []*websocket.Conn{conn},
	}
}

func createTopic(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	name := getSequence()
	topic := NewTopic(conn)
	hub.Cables[name] = topic

	msg, err := json.Marshal(&struct {
		Topic string `json:"topic"`
	}{
		Topic: name,
	})
	if err != nil {
		log.Println(err)
	}

	topic.Broadcast(&Message{
		msg:  msg,
		conn: nil,
	})

	go readPump(topic, conn)
}

func joinTopic(hub *Hub, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("topic ID: %v\n", vars["id"])
	id := vars["id"]

	if topic := hub.Cables[id]; topic != nil {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
		}
		topic.Register(conn)
		go readPump(topic, conn)
	}
}

func readPump(topic *Cable, conn *websocket.Conn) {
	log.Printf("Topic: %p, conn: %p opened", topic, conn)
	defer func() {
		topic.Unregister(conn)
		if len(topic.conns) == 0 {
		}
		conn.Close()
		log.Printf("Topic: %p, conn: %p closed", topic, conn)
	}()
	for {
		t, message, err := conn.ReadMessage()
		log.Println(t)
		if err != nil {
			log.Println(err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		log.Println(string(message))
		topic.Broadcast(&Message{
			msg:  message,
			conn: conn,
		})
	}
}
