package libfgg

import (
	"context"
	"os"
)

type Fgg struct {
	Server string
	Action string
}

func (f *Fgg) Topic() string {
	return f.Server + "/topic/"
}

func (f *Fgg) Send(ctx context.Context, list []string) {
	if len(list) == 0 {
		panic("Need File")
	}

	ws := NewWebSocketConn()
	ws.Start(ctx, f.Topic())
	go ws.Run(ctx)

	file, err := os.Open(list[0])
	if err != nil {
		panic(err)
	}
	transfer := &Transfer{
		Conn: ws,
		File: file,
	}
	transfer.Send()
	transfer.Run()
}

func (f *Fgg) Recv(ctx context.Context, list []string) {
	var file *os.File
	var err error
	if len(list) != 0 {
		file, err = os.Create(list[0])
		if err != nil {
			panic(err)
		}
	}

	ws := NewWebSocketConn()
	ws.Start(ctx, f.Topic())
	go ws.Run(ctx)

	transfer := &Transfer{
		Conn: ws,
		File: file,
	}
	transfer.Recv()
	transfer.Run()
}
