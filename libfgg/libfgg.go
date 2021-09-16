package libfgg

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"filegogo/libfgg/transfer"
	"filegogo/libfgg/webrtc"
	"filegogo/libfgg/websocket"

	pion "github.com/pion/webrtc/v3"
	"github.com/sb-im/jsonrpc-lite"
	log "github.com/sirupsen/logrus"
)

type Fgg struct {
	Tran *transfer.Transfer
	Conn Conn

	ws  *websocket.Conn
	rtc *webrtc.Conn

	errors chan error
	finish bool
	sender bool

	// Callbacks
	OnPreTran  func(*transfer.MetaFile)
	OnPostTran func(*transfer.MetaHash)
}

func NewFgg() *Fgg {
	return &Fgg{
		Tran:       transfer.NewTransfer(),
		errors:     make(chan error),
		OnPreTran:  func(meta *transfer.MetaFile) {},
		OnPostTran: func(meta *transfer.MetaHash) {},
	}
}

func (t *Fgg) Send(files []string) error {
	if err := t.Tran.Send(files); err != nil {
		return err
	}
	t.sender = true
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
		log.Warn(err)
		return err
	}
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

func (t *Fgg) Run() error {
	return <-t.errors
}

func (t *Fgg) onPreTran(meta *transfer.MetaFile) {
	t.OnPreTran(meta)
}

func (t *Fgg) onPostTran(meta *transfer.MetaHash) {
	t.OnPostTran(meta)
}

func (t *Fgg) GetFile() {
	data, _ := jsonrpc.NewNotify("getfile", nil).ToJSON()
	t.send(data, TypeStr)
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
		case "getfile":
			t.sendData()
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

			t.onPreTran(meta)
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
	t.send(data, TypeStr)
}

func (t *Fgg) reqdata() {
	data, _ := jsonrpc.NewNotify("reqdata", nil).ToJSON()
	t.send(data, TypeStr)
}

func (t *Fgg) reslist() {
	if t.sender {
		meta := t.Tran.GetMetaFile()
		data, _ := jsonrpc.NewNotify("filelist", meta).ToJSON()

		t.onPreTran(meta)
		t.send(data, TypeStr)
	}
}

func (t *Fgg) sendData() {
	data, err := t.Tran.Read()
	if err != nil {
		return
	}
	t.send(data, TypeBin)
}

func (t *Fgg) reqsum() {
	data, _ := jsonrpc.NewNotify("reqsum", nil).ToJSON()
	t.send(data, TypeStr)
}

func (t *Fgg) ressum() {
	meta := t.Tran.GetMetaHash()
	data, err := jsonrpc.NewNotify("ressum", meta).ToJSON()
	t.send(data, TypeStr)
	t.onPostTran(meta)

	// Need Wait websocket send data
	time.Sleep(time.Millisecond)
	t.Close(err)
}

func (t *Fgg) Verify(meta *transfer.MetaHash) {
	if t.Tran.VerifyHash(meta) {
		t.Close(nil)
	} else {
		t.Close(errors.New("md5 VerifyHash failed"))
	}
	t.onPostTran(meta)
}

func (t *Fgg) Close(err error) {
	//t.ws.Close()
	//t.rtc.Close()
	t.errors <- err
}
