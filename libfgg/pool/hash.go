package pool

import (
	"crypto/md5"
	"fmt"
	"hash"
)

type fileHash struct {
	hash   hash.Hash
	offset int64
}

func newFileHash() *fileHash {
	return &fileHash{
		hash: md5.New(),
	}
}

func (f *fileHash) onData(c *DataChunk, data []byte) error {
	if f.offset == c.Offset {
		n, err := f.hash.Write(data)
		f.offset += int64(n)
		return err
	}
	return nil
}

func (f *fileHash) sum() string {
	return fmt.Sprintf("%x", f.hash.Sum(nil))
}
