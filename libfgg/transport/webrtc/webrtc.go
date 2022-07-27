package webrtc

import (
	"context"
	"time"

	"filegogo/libfgg/transport/protocol"

	"github.com/pion/datachannel"
	log "github.com/sirupsen/logrus"
)

type Conn struct {
	conn      datachannel.ReadWriteCloser
	onMessage func([]byte, []byte)
}

func New(conn datachannel.ReadWriteCloser) *Conn {
	return &Conn{
		conn:      conn,
		onMessage: func([]byte, []byte) {},
	}
}

func (c *Conn) SetOnRecv(fn func(head, body []byte)) {
	c.onMessage = fn
}

func (c *Conn) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			c.conn.Close()
			return
		default:
		}
		data := make([]byte, 1024*64)
		n, err := c.conn.Read(data)
		if err != nil {
			time.Sleep(time.Millisecond)
			continue
		}
		log.Tracef("WebRTC DataChannel RECV count(%d): '%s'", n, data[:n])

		c.onMessage(protocol.Decode(data[:n]))
	}
}

func (c *Conn) Send(head, body []byte) error {
	data := protocol.Encode(head, body)

	n, err := c.conn.Write(data)
	log.Tracef("WebRTC DataChannel SEND count(%d): %s\n", n, data[:n])
	return err
}
