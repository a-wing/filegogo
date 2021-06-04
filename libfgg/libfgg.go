package libfgg

import (
	"encoding/json"
	"hash"
	"log"
	"os"
	"time"

	"github.com/SB-IM/jsonrpc-lite"
)

type Fgg struct {
	File *os.File
	Tran *Transfer
	Conn Conn
	Hash hash.Hash
	send bool
	run  bool

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
	for t.run {
		messageType, data, err := t.Conn.Recv()
		if err != nil {
			log.Fatal(err)
		}
		if messageType == TextMessage {
			rpc := jsonrpc.ParseObject(data)
			switch rpc.Method {
			case "reqlist":
				t.reslist()
			case "reqsdp":
				rtc := NewWebrtcConn()
				sign := make(chan bool)
				rtc.OnOpen = func() {
					sign <- true
				}
				rtc.RunAnswer(t.Conn)
				<-sign
				log.Println("WebRTC Connected")
				t.Conn = rtc

			case "reqdata":
				t.sendData()
			case "reqsum":
				t.ressum()
			case "ressum":
				hash := &MetaHash{}
				json.Unmarshal(*rpc.Params, hash)
				t.Verify(hash)
			case "filelist":
				t.reqsdp()

				rtc := NewWebrtcConn()
				sign := make(chan bool)
				rtc.OnOpen = func() {
					sign <- true
				}
				rtc.RunOffer(t.Conn)

				meta := &MetaFile{}
				json.Unmarshal(*rpc.Params, meta)

				t.Tran.SetMetaFile(meta)
				t.OnPreTran(meta)

				<-sign
				time.Sleep(time.Second)
				log.Println("WebRTC Connected")
				t.Conn = rtc

				t.reqdata()
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

func (t *Fgg) reqsdp() {
	data, _ := jsonrpc.NewNotify("reqsdp", nil).ToJSON()
	t.Conn.Send(TextMessage, data)
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
	t.Conn.Close()
	t.run = false
}

func (t *Fgg) Verify(meta *MetaHash) {
	log.Println()
	if t.Tran.VerifyHash(meta) {
		log.Println("md5 sum (ok): ", meta.Hash)
	} else {
		log.Println("source file ms5: ", meta.Hash)
	}

	t.run = false
}
