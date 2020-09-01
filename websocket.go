package main

import (
	"bytes"
	"log"
	"net/http"

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

func readPump(hub *Hub, id string, conn *websocket.Conn) {
	defer func() {
		log.Println("EEEEEEEEEEEEEEEe")
		//c.hub.unregister <- c
		//c.conn.Close()
		conn.Close()
	}()
	//c.conn.SetReadLimit(maxMessageSize)
	//c.conn.SetReadDeadline(time.Now().Add(pongWait))
	//c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
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
		msg := &Message{
			msg:  message,
			conn: conn,
		}
		hub.Broadcast(id, msg)
		//c.hub.broadcast <- message
	}
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	log.Println(r.URL)
	log.Println(r.Header)
	if err != nil {
		log.Println(err)
		return
	}

	vars := mux.Vars(r)
	log.Printf("ID: %v\n", vars["id"])

	id := vars["id"]
	if hub.Cables[id] == nil {
		cable := &Cable{
			Message: &Message{
				msg:  []byte("room id : " + id),
				conn: conn,
			},
			conns: []*websocket.Conn{conn},
		}
		hub.Cables[id] = cable
	} else {
		hub.Cables[id].Register(conn)
	}

	go readPump(hub, vars["id"], conn)
}

