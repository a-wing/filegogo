package socket

import (
	"context"
	"encoding/binary"
	"io"
	"net"

	"filegogo/libfgg/transport/protocol"

	log "github.com/sirupsen/logrus"
)

const (
	l1Length = protocol.LengthHead
	l2Length = protocol.LengthBody
)

type Conn struct {
	conn      net.Conn
	onMessage func([]byte, []byte)
}

func New(conn net.Conn) *Conn {
	return &Conn{
		conn:      conn,
		onMessage: func([]byte, []byte) {},
	}
}

func (c *Conn) Run(ctx context.Context) {
	l1l2 := make([]byte, l1Length+l2Length)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if _, err := io.ReadFull(c.conn, l1l2); err != nil {
				log.Error(err)
				return
			}

			l1, l2 := binary.BigEndian.Uint16(l1l2[:l1Length]), binary.BigEndian.Uint16(l1l2[l1Length:l1Length+l2Length])

			log.Debug(l1l2, l1, l2)

			head := make([]byte, l1)
			if _, err := io.ReadFull(c.conn, head); err != nil {
				log.Error(err)
				return
			}

			body := make([]byte, l2)
			if _, err := io.ReadFull(c.conn, body); err != nil {
				log.Error(err)
				return
			}
			c.onMessage(head, body)
		}
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
	if _, err := c.conn.Write(l1l2); err != nil {
		return err
	}
	if _, err := c.conn.Write(head); err != nil {
		return err
	}
	if _, err := c.conn.Write(body); err != nil {
		return err
	}
	return nil
}

func (c *Conn) SetOnRecv(fn func(head, body []byte)) {
	c.onMessage = fn
}
