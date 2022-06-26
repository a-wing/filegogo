package socket

import (
	"context"
	"encoding/binary"
	"io"
	"net"

	log "github.com/sirupsen/logrus"
)

type Conn struct {
	conn      net.Conn
	OnMessage func([]byte, []byte)
}

func New(conn net.Conn) *Conn {
	return &Conn{
		conn:      conn,
		OnMessage: func([]byte, []byte) {},
	}
}

func (c *Conn) Run(ctx context.Context) {
	l1l2 := make([]byte, 4)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if _, err := io.ReadFull(c.conn, l1l2); err != nil {
				log.Error(err)
				return
			}

			l1, l2 := binary.BigEndian.Uint16(l1l2[:2]), binary.BigEndian.Uint16(l1l2[2:4])

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
			c.OnMessage(head, body)
		}
	}
}

func (c *Conn) Send(head, body []byte) error {
	l1, l2 := len(head), len(body)

	l1b := make([]byte, 2)
	l2b := make([]byte, 2)
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
	c.OnMessage = fn
}
