package websocket

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var (
	Type2Bool = map[int]bool{
		1: true,
		2: false,
	}

	Bool2Type = map[bool]int{
		true:  1,
		false: 2,
	}
)

type Conn struct {
	Conn   *websocket.Conn
	token  string
	server string

	OnOpen    func()
	OnClose   func()
	OnError   func(error)
	OnMessage func([]byte, bool)
}

func NewConn(addr string) *Conn {
	return &Conn{
		server:    addr,
		OnOpen:    func() {},
		OnClose:   func() {},
		OnError:   func(error) {},
		OnMessage: func([]byte, bool) {},
	}
}

func (c *Conn) Connect() (err error) {
	c.Conn, _, err = websocket.DefaultDialer.Dial(c.server, nil)
	return
}

func (c *Conn) Close() error {
	if err := c.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
		log.Println("write close:", err)
		return err
	}
	return c.Conn.Close()
}

func (c *Conn) Run() {
	for {
		typ, data, err := c.Conn.ReadMessage()
		log.Tracef("WebSocket RECV %d: %s\n", typ, data)
		if err != nil {
			log.Warn(err)
			c.OnError(err)
		}
		c.OnMessage(data, Type2Bool[typ])
	}
}

func (c *Conn) Send(data []byte, typ bool) error {
	log.Tracef("WebSocket SEND %d: %s\n", Bool2Type[typ], data)
	return c.Conn.WriteMessage(Bool2Type[typ], data)
}
