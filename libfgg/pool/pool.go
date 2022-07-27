package pool

import (
	"errors"
	"fmt"
	"os"
	"sync"
)

type Meta struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Size int64  `json:"size"`
}

type Hash struct {
	File string `json:"name"`
	Hash string `json:"hash"`
}

type Pool struct {
	sender *os.File
	recver *os.File

	fileHash *fileHash

	meta *Meta
	hash *Hash

	doneCount int
	nextCount int

	mu sync.Mutex

	OnFinish   func()
	OnProgress func(c int64)

	chunkSize int64

	currentSize int64
	pendingSize int64
}

func New() *Pool {
	return &Pool{
		fileHash:   newFileHash(),
		OnFinish:   func() {},
		OnProgress: func(c int64) {},
		chunkSize:  32 * 1024,
	}
}

func (p *Pool) SetSend(file string) (err error) {
	p.sender, err = os.Open(file)
	return
}

func (p *Pool) SetRecv(file string) (err error) {
	p.recver, err = os.Create(file)
	return
}

func (p *Pool) SendMeta() (*Meta, error) {
	if p.meta != nil {
		return p.meta, nil
	}

	if file := p.sender; file != nil {
		// TODO: Maybe use https://github.com/gabriel-vasile/mimetype
		mimeType, _ := getFileContentType(file)
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types
		// mimeType := "application/octet-stream"

		stat, err := file.Stat()
		meta := &Meta{
			Name: file.Name(),
			Type: mimeType,
			Size: stat.Size(),
		}
		p.meta = meta
		return meta, err
	}

	return nil, errors.New("Sender not found file")
}

func (p *Pool) RecvMeta(meta *Meta) error {
	var err error
	if p.recver == nil {
		p.recver, err = os.Create(meta.Name)
	}
	p.meta = meta
	return err
}

func (p *Pool) SendHash() (*Hash, error) {
	return &Hash{
			File: p.sender.Name(),
			Hash: p.fileHash.sum(),
		}, func() error {
			if p.fileHash.offset != p.meta.Size {
				return errors.New("Not sum")
			}
			return nil
		}()
}

func (p *Pool) RecvHash(meta *Hash) error {
	if meta.Hash == p.fileHash.sum() {
		return nil
	}
	return fmt.Errorf("'%s' Not match '%s'", meta.Hash, p.fileHash.sum())
}
