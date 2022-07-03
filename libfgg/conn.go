package libfgg

import (
	"context"

	wsConn "filegogo/libfgg/transport/websocket"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

func (t *Fgg) UseWebsocket(addr string) error {
	log.Debug("websocket connect: ", addr)

	ws, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return err
	}

	conn := wsConn.New(ws)
	t.AddConn(conn)

	go conn.Run(context.Background())
	return nil
}
