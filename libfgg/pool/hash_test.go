package pool

import (
	"crypto/md5"
	"fmt"
	"testing"
)

func TestHash(t *testing.T) {
	data1 := []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	data2 := []byte("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	data3 := []byte("cccccccccccccccccccccccccccc")

	sum := md5.New()
	sum.Write(data1)

	fh := newFileHash()
	fh.onData(&DataChunk{
		Offset: 0,
		Length: int64(len(data1)),
	}, data1)

	if fh.sum() != fmt.Sprintf("%x", sum.Sum(nil)) {
		t.Error(fh.sum())
	}

	sum.Write(data2)

	fh.onData(&DataChunk{
		Offset: int64(len(data1)),
		Length: int64(len(data2)),
	}, data2)

	if fh.sum() != fmt.Sprintf("%x", sum.Sum(nil)) {
		t.Error(fh.sum())
	}

	// duplicate data2
	fh.onData(&DataChunk{
		Offset: int64(len(data1)),
		Length: int64(len(data2)),
	}, data2)

	if fh.sum() != fmt.Sprintf("%x", sum.Sum(nil)) {
		t.Error(fh.sum())
	}

	// duplicate data1
	fh.onData(&DataChunk{
		Offset: 0,
		Length: int64(len(data1)),
	}, data2)

	if fh.sum() != fmt.Sprintf("%x", sum.Sum(nil)) {
		t.Error(fh.sum())
	}

	sum.Write(data3)

	fh.onData(&DataChunk{
		Offset: int64(len(data1) + len(data2)),
		Length: int64(len(data3)),
	}, data3)

	if fh.sum() != fmt.Sprintf("%x", sum.Sum(nil)) {
		t.Error(fh.sum())
	}
}
