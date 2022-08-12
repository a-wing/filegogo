package websocket

import (
	"context"
	"sync"
	"time"

	"filegogo/libfgg/transport/protocol"

	"github.com/gorilla/websocket"

	log "github.com/sirupsen/logrus"
)

type Conn struct {
	conn      *websocket.Conn
	mutex     sync.Mutex
	onMessage func([]byte, []byte)
}

func New(conn *websocket.Conn) *Conn {
	return &Conn{
		conn:      conn,
		onMessage: func([]byte, []byte) {},
	}
}

func (c *Conn) SetOnRecv(fn func(head, body []byte)) {
	c.onMessage = fn
}

func (c *Conn) Run(ctx context.Context) {
	go func() {
		select {
		case <-ctx.Done():
			c.conn.Close()
		}
	}()

	for {
		typ, data, err := c.conn.ReadMessage()
		log.Tracef("WebSocket RECV %d: %s\n", typ, data)
		if err != nil {
			log.Error(err)
		}

		if data == nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		c.onMessage(protocol.Decode(data))
	}
}

func (c *Conn) Send(head, body []byte) error {
	data := protocol.Encode(head, body)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.conn.WriteMessage(websocket.BinaryMessage, data)
}
