package libfgg

import (
	"context"
	"encoding/json"
	"hash"
	"log"
	"os"
	"time"

	"filegogo/libfgg/webrtc"

	"github.com/SB-IM/jsonrpc-lite"
	pion "github.com/pion/webrtc/v3"
	"github.com/spf13/viper"
)

type Fgg struct {
	File *os.File
	Tran *Transfer
	Conn Conn
	Hash hash.Hash
	send bool
	run  bool

	ws  Conn
	rtc *webrtc.Conn

	cancel context.CancelFunc

	finish bool

	// Callbacks
	OnPreTran  func(*MetaFile)
	OnPostTran func(*MetaHash)
}

func (t *Fgg) Send() {
	t.Tran = NewTransfer(t.File)
	t.send = true
	t.run = true
	t.reslist()
}

func (t *Fgg) Recv() {
	t.Tran = NewTransfer(t.File)
	t.Tran.OnFinish = func() {
		t.finish = true
	}
	t.run = true
	t.reqlist()
}

func (t *Fgg) Run() {
	iceservers := &pion.Configuration{}
	viper.Unmarshal(iceservers)
	dd, _ := json.Marshal(iceservers)
	log.Println(string(dd))
	t.rtc = webrtc.NewConn(iceservers)

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

	t.ws = t.Conn
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
				hash := &MetaHash{}
				json.Unmarshal(*rpc.Params, hash)
				t.Verify(hash)
			case "filelist":
				t.rtc.Start()

				meta := &MetaFile{}
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

func (t *Fgg) Verify(meta *MetaHash) {
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
