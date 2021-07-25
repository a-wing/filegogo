package websocket

import (
	"encoding/json"
	"errors"
	"net/url"
	"path"
	"strings"

	"filegogo/server"

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

func fillAddr(addr string) string {
	if u, err := url.Parse(addr); err != nil {
		return ""
	} else {
		name := ""
		if arr := strings.Split(u.Path, "/"); len(arr) > 2 {
			name = arr[2]
		}
		u.Path = ""
		return u.String() + "/" + server.PrefixShare + "/" + name
	}
}

func NewConn(addr string) *Conn {
	return &Conn{
		server:    fillAddr(addr),
		OnOpen:    func() {},
		OnClose:   func() {},
		OnError:   func(error) {},
		OnMessage: func([]byte, bool) {},
	}
}

func (c *Conn) Server() string {
	return c.server
}

func (c *Conn) Connect() (err error) {
	c.Conn, _, err = websocket.DefaultDialer.Dial(c.authServer(), nil)
	if err != nil {
		log.Println(c.server)
		log.Println(c.authServer())
		log.Fatal("dial:", err)
		return
	}
	typ, data, err := c.Conn.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		return
	}
	log.Debugf("recv: %s", data)

	if typ != websocket.TextMessage {
		err = errors.New("Must is text message")
	}

	msg := &server.MessageHello{}
	if err = json.Unmarshal(data, msg); err != nil {
		return
	}
	c.updateServer(msg.Share)
	c.token = msg.Token
	return
}

func (c *Conn) authServer() string {
	if c.token == "" {
		return c.server
	} else {
		return c.server + "?token=" + c.token
	}
}

func (c *Conn) updateServer(share string) error {
	if u, err := url.Parse(c.server); err != nil {
		return err
	} else {
		u.Path = path.Join(server.PrefixShare, share)
		c.server = u.String()
		return nil
	}
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
