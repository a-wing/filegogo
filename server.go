package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type Message struct {
	conn *websocket.Conn
	msg  []byte
}

type Cable struct {
	Message *Message
	conns   []*websocket.Conn
}

func newCable(msg []byte, conn *websocket.Conn) *Cable {
	return &Cable{
		conns: []*websocket.Conn{conn},
		Message: &Message{
			conn: conn,
			msg:  msg,
		},
	}
}

func (this *Cable) Register(conn *websocket.Conn) {
	this.conns = append(this.conns, conn)
}

func (this *Cable) Unregister(conn *websocket.Conn) {
	for index, link := range this.conns {
		if link == conn {

			// Order is not important
			this.conns[index] = this.conns[len(this.conns)-1]
			this.conns = this.conns[:len(this.conns)-1]
		}
	}
}

func (this *Cable) Broadcast(msg *Message) {
	log.Println("run Broadcast")
	for _, conn := range this.conns {
		log.Println(&conn)
		if msg.conn != conn {
			err := conn.WriteMessage(websocket.TextMessage, msg.msg)
			if err != nil {
				log.Println(err)
				this.Unregister(conn)
			}

		}
	}
}

type Hub struct {
	Cables map[string]*Cable
}

func newHub() *Hub {
	return &Hub{
		//broadcast:  make(chan []byte),
		Cables: make(map[string]*Cable),
	}
}

func (this *Hub) Add(name string, cable *Cable) {
	this.Cables[name] = cable
}

func (this *Hub) Remove(name string) {
	delete(this.Cables, name)
}

func (this *Hub) Broadcast(name string, msg *Message) {
	if cable := this.Cables[name]; cable != nil {
		log.Println(name, "cable is: ", cable)
		cable.Broadcast(msg)
	}
}

