package libfgg

import (
	"context"
	"encoding/json"
	"time"

	"filegogo/libfgg/transfer"
	"filegogo/libfgg/webrtc"
	"filegogo/libfgg/websocket"

	"github.com/SB-IM/jsonrpc-lite"
	pion "github.com/pion/webrtc/v3"
	log "github.com/sirupsen/logrus"
)

type Fgg struct {
	Tran *transfer.Transfer
	Conn Conn

	ws  *websocket.Conn
	rtc *webrtc.Conn

	finish bool
	cancel context.CancelFunc

	// Callbacks
	OnShare    func(addr string)
	OnPreTran  func(*transfer.MetaFile)
	OnPostTran func(*transfer.MetaHash)
}

func NewFgg() *Fgg {
	return &Fgg{
		Tran:       transfer.NewTransfer(),
		OnShare:    func(addr string) {},
		OnPreTran:  func(meta *transfer.MetaFile) {},
		OnPostTran: func(meta *transfer.MetaHash) {},
	}
}

func (t *Fgg) Send(files []string) error {
	if err := t.Tran.Send(files); err != nil {
		return err
	}
	t.reslist()
	return nil
}

func (t *Fgg) Recv(files []string) error {
	if err := t.Tran.Recv(files); err != nil {
		return err
	}
	t.Tran.OnFinish = func() {
		t.finish = true
	}
	t.reqlist()
	return nil
}

func (t *Fgg) UseWebsocket(addr string) error {
	log.Debug("websocket connect: ", addr)
	t.ws = websocket.NewConn(addr)
	t.ws.OnMessage = t.recv
	if err := t.ws.Connect(); err != nil {
		log.Trace(t.ws.Server())
		log.Warn(err)
		return err
	}
	t.OnShare(t.ws.Server())
	t.Conn = t.ws

	go t.ws.Run()
	return nil
}

func (t *Fgg) UseWebRTC(IceServers *pion.Configuration) {
	t.rtc = webrtc.NewConn(IceServers)
	t.rtc.OnMessage = t.recv

	t.rtc.OnSignSend = func(data []byte) {
		rpc := jsonrpc.NewNotify("webrtc", nil)
		RawMessage := json.RawMessage(data)
		rpc.Params = &RawMessage
		raw, _ := rpc.ToJSON()

		t.send(raw, true)
	}

	t.rtc.OnOpen = func() {
		log.Println("WebRTC Connected")

		// Debug use Disable Websocket OnMessage
		// t.ws.OnMessage = func(b1 []byte, b2 bool) {}

		t.Conn = t.rtc
		go t.rtc.Run()
	}
}

func (t *Fgg) RunWebRTC() {
	ctx, cancel := context.WithCancel(context.Background())
	onOpen := t.rtc.OnOpen
	t.rtc.OnOpen = func() {
		onOpen()
		cancel()
	}
	t.rtc.Start()
	ticker := time.NewTicker(3 * time.Second)
	select {
	case <-ticker.C:
		log.Warn("WebRTC timeout")
	case <-ctx.Done():
	}

	return
}

func (t *Fgg) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	t.cancel = cancel
	<-ctx.Done()
}

func (t *Fgg) GetFile() {
	t.reqdata()
}

func (t *Fgg) send(data []byte, typ bool) error {
	return t.Conn.Send(data, typ)
}

func (t *Fgg) recv(data []byte, typ bool) {
	if typ {
		rpc := jsonrpc.ParseObject(data)
		switch rpc.Method {
		case "webrtc":
			t.rtc.SignRecv(*rpc.Params)
		case "reqlist":
			t.reslist()
		case "reqdata":
			t.sendData()
		case "reqsum":
			t.ressum()
		case "ressum":
			hash := &transfer.MetaHash{}
			json.Unmarshal(*rpc.Params, hash)
			t.Verify(hash)
		case "filelist":
			meta := &transfer.MetaFile{}
			json.Unmarshal(*rpc.Params, meta)

			t.Tran.SetMetaFile(meta)

			// The hook maybe block
			go t.OnPreTran(meta)
		}
	} else {
		t.Tran.Write(data)
		if t.finish {
			t.reqsum()
		} else {
			t.reqdata()
		}
	}
}

func (t *Fgg) reqlist() {
	data, _ := jsonrpc.NewNotify("reqlist", nil).ToJSON()
	t.send(data, true)
}

func (t *Fgg) reqdata() {
	data, _ := jsonrpc.NewNotify("reqdata", nil).ToJSON()
	t.send(data, true)
}

func (t *Fgg) reslist() {
	meta := t.Tran.GetMetaFile()
	data, _ := jsonrpc.NewNotify("filelist", meta).ToJSON()

	t.OnPreTran(meta)
	t.send(data, true)
}

func (t *Fgg) sendData() {
	data, err := t.Tran.Read()
	if err != nil {
		return
	}
	//t.Conn.Send(BinaryMessage, data)
	t.send(data, false)
}

func (t *Fgg) reqsum() {
	data, _ := jsonrpc.NewNotify("reqsum", nil).ToJSON()
	//t.Conn.Send(TextMessage, data)
	t.send(data, true)
}

func (t *Fgg) ressum() {
	meta := t.Tran.GetMetaHash()
	data, _ := jsonrpc.NewNotify("ressum", meta).ToJSON()
	//t.Conn.Send(TextMessage, data)
	t.send(data, true)

	// Need Wait websocket send data
	time.Sleep(time.Second)
	t.Close()
}

func (t *Fgg) Verify(meta *transfer.MetaHash) {
	log.Println()
	if t.Tran.VerifyHash(meta) {
		log.Println("md5 sum (ok): ", meta.Hash)
	} else {
		log.Println("source file ms5: ", meta.Hash)
	}
	t.Close()
}

func (t *Fgg) Close() {
	//t.ws.Close()
	//t.rtc.Close()
	t.cancel()
}
