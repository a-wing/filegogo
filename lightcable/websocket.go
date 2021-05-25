package lightcable

import (
	"log"

	"github.com/gorilla/websocket"
)

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
