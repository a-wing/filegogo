package libfgg

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"filegogo/lightcable"

	"github.com/SB-IM/jsonrpc-lite"
	"github.com/gorilla/websocket"
)

type WebSocketConn struct {
	Conn           *websocket.Conn
	Token          string
	Server         string
	OnConnected    func()
	OnDisconnected func()
}

func NewWebSocketConn() *WebSocketConn {
	return &WebSocketConn{}
}

func (c *WebSocketConn) Send(t int, data []byte) error {
	return c.Conn.WriteMessage(t, data)
}

func (c *WebSocketConn) Recv() (int, []byte, error) {
	return c.Conn.ReadMessage()
}

func (ws *WebSocketConn) Start(ctx context.Context, addr string) {
	c, _, err := websocket.DefaultDialer.Dial(ShareToWebSocket(addr), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	ws.Conn = c
	messageType, message, err := c.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		return
	}
	log.Printf("recv: %s", message)

	switch messageType {
	case websocket.TextMessage:
		rpc, err := jsonrpc.Parse(message)
		if err != nil {
			log.Println("read:", err)
		} else {
			switch rpc.Method {
			case "server":
				msg := &lightcable.MessageHello{}
				if err := json.Unmarshal(*rpc.Params, msg); err == nil {
					topic := ShareToWebSocket(addr + msg.Topic)

					fmt.Println(topic)
					fmt.Println(WebSocketToShare(topic))
					fmt.Println("=========")

					ws.Server = topic
					ws.Token = msg.Token
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
