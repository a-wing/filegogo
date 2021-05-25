package lightcable

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	conn *websocket.Conn
	code int
	data []byte
}

type Topic struct {
	Name  string
	mutex sync.Mutex
	conns []*websocket.Conn
}

func NewTopic(name string, conn *websocket.Conn) *Topic {
	return &Topic{
		Name:  name,
		conns: []*websocket.Conn{conn},
	}
}

func (this *Topic) Register(conn *websocket.Conn) {
	this.mutex.Lock()
	this.conns = append(this.conns, conn)
	this.mutex.Unlock()
}

func (this *Topic) Unregister(conn *websocket.Conn) {
	this.mutex.Lock()
	for index, link := range this.conns {
		if link == conn {

			// Order is not important
			this.conns[index] = this.conns[len(this.conns)-1]
			this.conns = this.conns[:len(this.conns)-1]
		}
	}
	this.mutex.Unlock()
}

func (this *Topic) Broadcast(msg *Message) {
	this.mutex.Lock()
	for _, conn := range this.conns {
		if msg.conn != conn {
			if err := conn.WriteMessage(msg.code, msg.data); err != nil {
				log.Println(err)
				this.Unregister(conn)
			}

		}
	}
	this.mutex.Unlock()
}
