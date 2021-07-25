package lightcable

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	topic map[string]*topic

	// Register requests from the clients.
	register chan *Client

	// Inbound messages from the clients.
	broadcast chan []byte

	// Unregister requests from clients.
	unregister chan *Client

	// TODO:
	// onConnClose(*Client, error)
	// OnMessage(*message)
	//onConnect func(*Client)
	// hook onConnect
	// hook disconnected
	// hook ommessage
	// func send message
}

func NewServer() *Server {
	return &Server{
		topic: make(map[string]*topic),

		register:   make(chan *Client),
		broadcast:  make(chan []byte),
		unregister: make(chan *Client),
	}
}

func (s *Server) Run(ctx context.Context) {
	for {
		select {
		// unregister must first
		// close and open concurrency
		case c := <-s.unregister:
			delete(s.topic, c.cable)
		case c := <-s.register:
			c.topic = s.topic[c.cable]
			if c.topic == nil {
				c.topic = NewTopic(c.cable, s)
				go c.topic.run(ctx)
				s.topic[c.cable] = c.topic
			}
			c.topic.register <- c

		case <-ctx.Done():
			// safe Close
		}
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
	}
	token := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", &conn)))
	s.JoinCable(r.URL.Path, token, conn)
}

func (s *Server) JoinCable(cable, label string, conn *websocket.Conn) error {
	select {
	case s.register <- &Client{
		cable: cable,
		label: label,
		conn:  conn,
		send:  make(chan message, 256),
	}:
		return nil
	default:
		return errors.New("join failure")
	}
}
