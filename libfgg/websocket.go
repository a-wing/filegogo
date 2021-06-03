package libfgg

import (
	"context"
	"encoding/json"
	"net/url"

	log "github.com/sirupsen/logrus"

	"filegogo/lightcable"

	"github.com/SB-IM/jsonrpc-lite"
	"github.com/gorilla/websocket"
)

type WebSocketConn struct {
	Conn    *websocket.Conn
	Token   string
	Server  string
	OnOpen  func()
	OnClose func()
}

func NewWebSocketConn() *WebSocketConn {
	return &WebSocketConn{
		OnOpen:  func() {},
		OnClose: func() {},
	}
}

func (c *WebSocketConn) SetServer(addr string, msg *lightcable.MessageHello) {
	if msg == nil {
		c.Server = addr
	} else {
		if u, err := url.Parse(addr); err != nil {
			return
		} else {
			u.Path = "/topic/" + msg.Topic
			c.Token = msg.Token
			c.Server = u.String()
		}
	}
}

func (c *WebSocketConn) Send(t int, data []byte) error {
	return c.Conn.WriteMessage(t, data)
}

func (c *WebSocketConn) Recv() (int, []byte, error) {
	return c.Conn.ReadMessage()
}

func (c *WebSocketConn) Close() error {
	return c.Conn.Close()
}

func (ws *WebSocketConn) Start(ctx context.Context, addr string) {
	c, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	ws.Conn = c
	messageType, message, err := c.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		return
	}
	log.Debugf("recv: %s", message)

	switch messageType {
	case websocket.TextMessage:
		rpc, err := jsonrpc.Parse(message)
		if err != nil {
			log.Fatalln("read:", err)
		} else {
			switch rpc.Method {
			case "server":
				msg := &lightcable.MessageHello{}
				if err := json.Unmarshal(*rpc.Params, msg); err == nil {
					ws.SetServer(addr, msg)
					ws.OnOpen()
				}
			default:
			}
		}

	case websocket.BinaryMessage:
	}
}

func (c *WebSocketConn) Run(ctx context.Context) {
	select {
	case <-ctx.Done():
		err := c.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Println("write close:", err)
			return
		}
	}
}
