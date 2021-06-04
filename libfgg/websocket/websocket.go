package websocket

import (
	"encoding/json"
	"errors"
	"net/url"
	"path"
	"strings"

	"filegogo/lightcable"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Conn struct {
	Conn   *websocket.Conn
	token  string
	server string
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
		return u.String() + "/" + lightcable.PrefixShare + "/" + name
	}
}

func NewConn(addr string) *Conn {
	return &Conn{
		server: fillAddr(addr),
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

	msg := &lightcable.MessageHello{}
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
		u.Path = path.Join(lightcable.PrefixShare, share)
		c.server = u.String()
		return nil
	}
}

func (c *Conn) Send(typ int, data []byte) error {
	return c.Conn.WriteMessage(typ, data)
}

func (c *Conn) Recv() (int, []byte, error) {
	return c.Conn.ReadMessage()
}

func (c *Conn) Close() error {
	if err := c.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
		log.Println("write close:", err)
		return err
	}
	return c.Conn.Close()
}
