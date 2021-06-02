package libfgg

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"hash"
	"io"
	"net/http"
	"os"
)

type MetaFile struct {
	File string `json:"file"`
	Type string `json:"type"`
	Size int64  `json:"size"`
}

type MetaHash struct {
	Hash string `json:"hash"`
}

type Transfer struct {
	File *os.File
	Hash hash.Hash

	metaFile *MetaFile
	//metaHash *MetaHash

	finish bool
	// progress total size
	count int64

	OnFinish   func()
	OnProgress func(c int64)

	// tmp buffer, because it mimetype
	buft      io.Reader
	chunkSize int
}

func NewTransfer(file *os.File) *Transfer {
	return &Transfer{
		File:       file,
		Hash:       md5.New(),
		OnFinish:   func() {},
		OnProgress: func(c int64) {},
		chunkSize:  1024,
	}
}

func (t *Transfer) SetMetaFile(meta *MetaFile) {
	if t.File == nil {
		t.File, _ = os.Create(meta.File)
	}
	t.metaFile = meta
}

func (t *Transfer) GetMetaFile() *MetaFile {
	// TODO: Maybe use https://github.com/gabriel-vasile/mimetype
	mimeType, _ := t.getFileContentType()
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types
	// mimeType := "application/octet-stream"

	stat, _ := t.File.Stat()
	meta := &MetaFile{
		File: t.File.Name(),
		Type: mimeType,
		Size: stat.Size(),
	}
	t.metaFile = meta
	return meta
}

func (t *Transfer) getFileContentType() (string, error) {
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	c, err := t.File.Read(buffer)
	if err != nil {
		return "", err
	}

	t.buft = bytes.NewReader(buffer[:c])

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer[:c])
	return contentType, nil
}

func (t *Transfer) getHash() string { return fmt.Sprintf("%x", t.Hash.Sum(nil)) }

func (t *Transfer) GetMetaHash() *MetaHash {
	return &MetaHash{
		Hash: t.getHash(),
	}
}

func (t *Transfer) VerifyHash(meta *MetaHash) bool {
	return meta.Hash == t.getHash()
}

func (t *Transfer) Read() ([]byte, error) {
	if t.finish {
		t.OnFinish()
		return []byte{}, nil
	}
	data := make([]byte, t.chunkSize)
	c, err := io.MultiReader(t.buft, t.File).Read(data)
	if err != nil {
		return data[:c], err
	}
	t.Hash.Write(data[:c])

	t.count += int64(c)
	t.OnProgress(int64(c))
	if t.count >= t.metaFile.Size {
		t.OnFinish()
		t.finish = true
		t.File.Close()
	}
	return data[:c], err
}

func (t *Transfer) Write(data []byte) error {
	if t.finish {
		t.OnFinish()
		return nil
	}
	c, err := t.File.Write(data)
	t.Hash.Write(data)
	t.count += int64(c)
	t.OnProgress(int64(c))

	if t.count >= t.metaFile.Size {
		t.OnFinish()
		t.finish = true
		t.File.Close()
	}
	return err
}
