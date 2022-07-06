package webrtc

import (
	"context"
	"encoding/binary"

	"github.com/pion/datachannel"
	log "github.com/sirupsen/logrus"
)

const (
	l1Length = 2
	l2Length = 2
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

		l1, l2 := binary.BigEndian.Uint16(data[:l1Length]), binary.BigEndian.Uint16(data[l1Length:l1Length+l2Length])

		payload := data[l1Length+l2Length:]
		head := payload[:l1]
		body := payload[l1 : l1+l2]
		c.onMessage(head, body)
	}
}

func (c *Conn) Send(head, body []byte) error {
	l1, l2 := len(head), len(body)

	l1b := make([]byte, l1Length)
	l2b := make([]byte, l2Length)
	binary.BigEndian.PutUint16(l1b, uint16(l1))
	binary.BigEndian.PutUint16(l2b, uint16(l2))
	l1l2 := append(l1b, l2b...)

	log.Debug(l1, l2, l1l2)

	data := append(append(l1l2, head...), body...)

	n, err := c.conn.Write(data)
	log.Tracef("WebRTC DataChannel SEND count(%d): %s\n", n, data[:n])
	return err
}
