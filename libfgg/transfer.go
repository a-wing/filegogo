package libfgg

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/SB-IM/jsonrpc-lite"
)

type Transfer struct {
	File *os.File
	Conn Conn
	Hash hash.Hash
	send bool
	buft io.Reader
	info FileList
	rate int64
	run  bool

	// Callbacks
	OnProgress func(c int64)
	OnPreTran  func(*FileList)
	OnPostTran func()
}

func (t *Transfer) Send() {
	t.send = true
	t.run = true
	t.reslist()
}

func (t *Transfer) Recv() {
	t.run = true
	t.reqlist()
}

func (t *Transfer) Run() {
	t.Hash = md5.New()
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
				rtc.RunAnswer(t.Conn)
				<-rtc.sign
				fmt.Println("WebRTC Connected")
				t.Conn = rtc

			case "reqdata":
				t.sendData()
			case "reqsum":
				t.ressum()
			case "ressum":
				sum := &Sum{}
				//log.Println(string(*rpc.Params))
				json.Unmarshal(*rpc.Params, sum)
				t.Verify(sum.CheckSum)
			case "filelist":
				t.reqsdp()

				rtc := NewWebrtcConn()
				rtc.RunOffer(t.Conn)

				list := &FileList{}
				json.Unmarshal(*rpc.Params, list)
				t.createFile(list)
				t.OnPreTran(list)

				<-rtc.sign
				time.Sleep(time.Second)
				fmt.Println("WebRTC Connected")
				t.Conn = rtc

				t.reqdata()
			}
		} else {
			t.File.Write(data)
			io.WriteString(t.Hash, string(data))
			t.rate += int64(len(data))
			t.OnProgress(int64(len(data)))
			if t.rate >= t.info.Size {
				t.File.Close()
				t.reqsum()
			} else {
				t.reqdata()
			}
		}
	}
}

func (t *Transfer) reqsdp() {
	data, _ := jsonrpc.NewNotify("reqsdp", nil).ToJSON()
	t.Conn.Send(TextMessage, data)
}

func (t *Transfer) reqlist() {
	data, _ := jsonrpc.NewNotify("reqlist", nil).ToJSON()
	t.Conn.Send(TextMessage, data)
}

func (t *Transfer) reqdata() {
	data, _ := jsonrpc.NewNotify("reqdata", nil).ToJSON()
	t.Conn.Send(TextMessage, data)
}

type FileList struct {
	File string `json:"file"`
	Type string `json:"type"`
	Size int64  `json:"size"`
}

func (t *Transfer) createFile(list *FileList) (err error) {
	if t.File == nil {
		t.File, err = os.Create(list.File)
	}
	t.info = *list
	return
}

func (t *Transfer) GetFileContentType(out *os.File) (string, error) {
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	c, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	t.buft = bytes.NewReader(buffer[:c])

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer[:c])
	return contentType, nil
}

func (t *Transfer) reslist() {

	// TODO: has use io.Reader here
	// Maybe use https://github.com/gabriel-vasile/mimetype
	mimeType, _ := t.GetFileContentType(t.File)
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types
	// mimeType := "application/octet-stream"

	stat, _ := t.File.Stat()

	fileinfo := FileList{
		File: t.File.Name(),
		Type: mimeType,
		Size: stat.Size(),
	}

	data, _ := jsonrpc.NewNotify("filelist", fileinfo).ToJSON()

	t.OnPreTran(&fileinfo)
	t.Conn.Send(TextMessage, data)
}

func (t *Transfer) sendData() {
	data := make([]byte, 1024)
	count, err := io.MultiReader(t.buft, t.File).Read(data)
	//log.Println(string(data))
	if err != nil {
		//log.Fatal(string(data) , err)
		return
	}
	io.WriteString(t.Hash, string(data[:count]))
	t.Conn.Send(BinaryMessage, data[:count])
	t.rate += int64(count)
	t.OnProgress(int64(count))
}

type Sum struct {
	CheckSum string `json:"checksum"`
}

func (t *Transfer) reqsum() {
	data, _ := jsonrpc.NewNotify("reqsum", nil).ToJSON()
	t.Conn.Send(TextMessage, data)
}

func (t *Transfer) ressum() {
	data, _ := jsonrpc.NewNotify("ressum", &Sum{
		CheckSum: t.getsum(),
	}).ToJSON()
	t.Conn.Send(TextMessage, data)

	// Need Wait websocket send data
	time.Sleep(time.Second)
	t.run = false
}

func (t *Transfer) Verify(sum string) {
	if t.getsum() == sum {
		log.Println("md5 sum (ok): ", sum)
	} else {
		log.Println("send ms5: ", sum)
		log.Println("recv ms5: ", t.getsum())
	}

	t.run = false
}

func (t *Transfer) getsum() string {
	return fmt.Sprintf("%x", t.Hash.Sum(nil))
}
