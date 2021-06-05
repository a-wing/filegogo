package transfer

import (
	"crypto/md5"
	"errors"
	"fmt"
	"hash"
	"net/http"
	"os"
)

type MetaFile struct {
	File string `json:"file"`
	Type string `json:"type"`
	Size int64  `json:"size"`
}

type MetaHash struct {
	File string `json:"file"`
	Hash string `json:"hash"`
}

type Transfer struct {
	file *os.File
	hash hash.Hash

	metaFile *MetaFile
	metaHash *MetaHash

	finish bool
	// progress total size
	count int64

	OnFinish   func()
	OnProgress func(c int64)

	chunkSize int
}

func NewTransfer() *Transfer {
	return &Transfer{
		hash:       md5.New(),
		OnFinish:   func() {},
		OnProgress: func(c int64) {},
		chunkSize:  1024,
	}
}

func (t *Transfer) Recv(files []string) (err error) {
	if len(files) != 0 {
		t.file, err = os.Create(files[0])
	}
	return
}

func (t *Transfer) Send(files []string) (err error) {
	if len(files) == 0 {
		return errors.New("Need File")
	}

	// TODO:
	// Need Support Multiple files
	t.file, err = os.Open(files[0])
	return
}

func (t *Transfer) SetMetaFile(meta *MetaFile) (err error) {
	if t.file == nil {
		t.file, err = os.Create(meta.File)
	}
	t.metaFile = meta
	return
}

func (t *Transfer) GetMetaFile() *MetaFile {
	// TODO: Maybe use https://github.com/gabriel-vasile/mimetype
	mimeType, _ := t.getFileContentType()
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types
	// mimeType := "application/octet-stream"

	stat, _ := t.file.Stat()
	meta := &MetaFile{
		File: t.file.Name(),
		Type: mimeType,
		Size: stat.Size(),
	}
	t.metaFile = meta
	return meta
}

func (t *Transfer) getFileContentType() (string, error) {
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	c, err := t.file.Read(buffer)
	defer t.file.Seek(0, 0)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer[:c])
	return contentType, nil
}

func (t *Transfer) getHash() string { return fmt.Sprintf("%x", t.hash.Sum(nil)) }

func (t *Transfer) GetMetaHash() *MetaHash {
	return &MetaHash{
		File: t.file.Name(),
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
	c, err := t.file.Read(data)
	if err != nil {
		return data[:c], err
	}
	t.hash.Write(data[:c])

	t.count += int64(c)
	t.OnProgress(int64(c))
	if t.count >= t.metaFile.Size {
		t.OnFinish()
		t.finish = true
		t.file.Close()
	}
	return data[:c], err
}

func (t *Transfer) Write(data []byte) error {
	if t.finish {
		t.OnFinish()
		return nil
	}
	c, err := t.file.Write(data)
	t.hash.Write(data)
	t.count += int64(c)
	t.OnProgress(int64(c))

	if t.count >= t.metaFile.Size {
		t.OnFinish()
		t.finish = true
		t.file.Close()
	}
	return err
}
