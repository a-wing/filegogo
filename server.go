package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type Message struct {
	conn *websocket.Conn
	msg  []byte
}

type Topic struct {
	Name  string
	conns []*websocket.Conn
}

func NewTopic(conn *websocket.Conn) *Topic {
	return &Topic{
		conns: []*websocket.Conn{conn},
	}
}

func (this *Topic) Register(conn *websocket.Conn) {
	this.conns = append(this.conns, conn)
}

func (this *Topic) Unregister(conn *websocket.Conn) {
	for index, link := range this.conns {
		if link == conn {

			// Order is not important
			this.conns[index] = this.conns[len(this.conns)-1]
			this.conns = this.conns[:len(this.conns)-1]
		}
	}
}

func (this *Topic) Broadcast(msg *Message) {
	for _, conn := range this.conns {
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
	Topics map[string]*Topic
}

func NewHub() *Hub {
	return &Hub{
		Topics: make(map[string]*Topic),
	}
}

func (this *Hub) Add(name string, topic *Topic) {
	this.Topics[name] = topic
}

func (this *Hub) Remove(name string) {
	delete(this.Topics, name)
}

func (this *Hub) Broadcast(name string, msg *Message) {
	if topic := this.Topics[name]; topic != nil {
		log.Println(name, "topic is: ", topic)
		topic.Broadcast(msg)
	}
}
