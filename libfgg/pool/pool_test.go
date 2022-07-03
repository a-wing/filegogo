package pool

import (
	"os"
	"path"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	// 100M
	sendPath, _, err := HelpCreateTmpFile("filegogo_libfgg_pool_send", 100*1024*1024)
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(sendPath)

	recvPath := path.Join(os.TempDir(), "filegogo_libfgg_pool_recv")
	defer os.Remove(recvPath)

	sender := New()
	if err := sender.SetSend(sendPath); err != nil {
		t.Error(err)
	}

	recver := New()
	if err := recver.SetRecv(recvPath); err != nil {
		t.Error(err)
	}

	meta, err := sender.SendMeta()
	if err != nil {
		t.Error(err)
	}

	if err := recver.RecvMeta(meta); err != nil {
		t.Error(err)
	}

	running := true
	recver.OnFinish = func() {
		running = false
	}
	for running {
		c := recver.Next()
		if c == nil {
			time.Sleep(100 * time.Millisecond)
		} else {
			data, err := sender.SendData(c)
			if err != nil {
				t.Error(err)
			}

			if err := recver.RecvData(c, data); err != nil {
				t.Error(err)
			}
		}
	}

	hash, err := sender.SendHash()
	if err != nil {
		t.Error(err)
	}

	if err := recver.RecvHash(hash); err != nil {
		t.Error(err)
	}
}
