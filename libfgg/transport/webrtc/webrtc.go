package webrtc

import (
	"context"

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
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		select {
		case <-ctx.Done():
			c.conn.Close()
		}
	}()

	for {
		// TODO:
		data := make([]byte, 1024*64)
		n, err := c.conn.Read(data)

		log.Tracef("WebRTC DataChannel RECV count(%d): %s\n", n, data[:n])
		if err != nil {
			log.Error(err)
			cancel()
			return
		}

		c.onMessage(protocol.Decode(data[:n]))
	}
}

func (c *Conn) Send(head, body []byte) error {
	data := protocol.Encode(head, body)

	n, err := c.conn.Write(data)
	log.Tracef("WebRTC DataChannel SEND count(%d): %s\n", n, data[:n])
	return err
}
