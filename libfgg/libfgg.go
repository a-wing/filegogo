package libfgg

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"filegogo/libfgg/transfer"
	"filegogo/libfgg/webrtc"
	"filegogo/libfgg/websocket"

	"github.com/SB-IM/jsonrpc-lite"
	pion "github.com/pion/webrtc/v3"
)

type Fgg struct {
	Tran *transfer.Transfer
	Conn Conn
	send bool
	run  bool

	ws  *websocket.Conn
	rtc *webrtc.Conn

	IceServers *pion.Configuration

	cancel context.CancelFunc

	finish bool

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
	t.send = true
	t.run = true
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
	t.run = true
	t.reqlist()
	return nil
}

func (t *Fgg) Start(addr string) {
	log.Println(addr)
	t.ws = websocket.NewConn(addr)
	if err := t.ws.Connect(); err != nil {
		log.Println(t.ws.Server())
		log.Fatal(err)
	}
	t.OnShare(t.ws.Server())
	t.Conn = t.ws
}

func (t *Fgg) Run() {
	// === WebRTC ===
	t.rtc = webrtc.NewConn(t.IceServers)

	t.rtc.OnSignSend = func(data []byte) {
		rpc := jsonrpc.NewNotify("webrtc", nil)
		RawMessage := json.RawMessage(data)
		rpc.Params = &RawMessage
		raw, _ := rpc.ToJSON()

		t.Conn.Send(TextMessage, raw)
	}

	t.rtc.OnOpen = func() {
		log.Println("WebRTC Connected")

		t.Conn = t.rtc
		go t.doRun()
		if !t.send {
			t.reqdata()
		}
	}

	//t.ws = t.Conn
	go t.doRun()
	ctx, cancel := context.WithCancel(context.Background())
	t.cancel = cancel
	<-ctx.Done()
}

func (t *Fgg) doRun() {
	for t.run {
		messageType, data, err := t.Conn.Recv()
		if err != nil {
			log.Fatal(err)
		}
		if messageType == TextMessage {
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
				t.rtc.Start()

				meta := &transfer.MetaFile{}
				json.Unmarshal(*rpc.Params, meta)

				t.Tran.SetMetaFile(meta)
				t.OnPreTran(meta)
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
}

func (t *Fgg) reqlist() {
	data, _ := jsonrpc.NewNotify("reqlist", nil).ToJSON()
	t.Conn.Send(TextMessage, data)
}

func (t *Fgg) reqdata() {
	data, _ := jsonrpc.NewNotify("reqdata", nil).ToJSON()
	t.Conn.Send(TextMessage, data)
}

func (t *Fgg) reslist() {
	meta := t.Tran.GetMetaFile()
	data, _ := jsonrpc.NewNotify("filelist", meta).ToJSON()

	t.OnPreTran(meta)
	t.Conn.Send(TextMessage, data)
}

func (t *Fgg) sendData() {
	data, err := t.Tran.Read()
	if err != nil {
		return
	}
	t.Conn.Send(BinaryMessage, data)
}

func (t *Fgg) reqsum() {
	data, _ := jsonrpc.NewNotify("reqsum", nil).ToJSON()
	t.Conn.Send(TextMessage, data)
}

func (t *Fgg) ressum() {
	meta := t.Tran.GetMetaHash()
	data, _ := jsonrpc.NewNotify("ressum", meta).ToJSON()
	t.Conn.Send(TextMessage, data)

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
	t.rtc.Close()
	t.run = false
	t.cancel()
}
