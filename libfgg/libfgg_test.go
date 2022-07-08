package libfgg

import (
	"context"
	"fmt"
	"net"
	"os"
	"path"
	"testing"
	"time"

	"filegogo/libfgg/pool"
	"filegogo/libfgg/transport/socket"
)

func TestLibFgg(t *testing.T) {
	// 100M
	sendPath, _, err := pool.HelpCreateTmpFile("filegogo_libfgg_send", 100*1024*1024)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(sendPath)

	recvPath := path.Join(os.TempDir(), "filegogo_libfgg_recv")
	defer os.Remove(recvPath)

	sendConn, recvConn := net.Pipe()
	sendFgg := NewFgg()
	recvFgg := NewFgg()

	sendSocket := socket.New(sendConn)
	recvSocket := socket.New(recvConn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go sendSocket.Run(ctx)
	go recvSocket.Run(ctx)

	sendFgg.AddConn(sendSocket)
	recvFgg.AddConn(recvSocket)

	sendFgg.SetSend(sendPath)
	recvFgg.SetRecv(recvPath)

	sendFgg.UseWebRTC(nil)
	recvFgg.UseWebRTC(nil)

	recvFgg.OnPreTran = func(c *pool.Meta) {
		fmt.Printf("%+v\n", c)
	}

	recvFgg.OnPostTran = func(c *pool.Hash) {
		fmt.Printf("%+v\n", c)
		cancel()
	}
	recvFgg.GetMeta()
	recvFgg.Run(ctx)
	fmt.Println("GGGGGGGGGGGGG")

	//sendFgg.Run()
	//recvFgg.Run()
}
