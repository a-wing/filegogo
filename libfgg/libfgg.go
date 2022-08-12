package libfgg

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"filegogo/libfgg/pool"
	"filegogo/libfgg/transport"

	"github.com/sb-im/jsonrpc-lite"
	log "github.com/sirupsen/logrus"
)

var uniqueID uint64

func getUniqueID() string {
	return strconv.FormatUint(atomic.AddUint64(&uniqueID, 1), 10)
}

var (
	loopWait        = 10 * time.Millisecond
	maxPendingCount = 100
)

const (
	methodMeta = "meta"
	methodData = "data"
	methodHash = "hash"
)

type call struct {
	req  *jsonrpc.Jsonrpc
	res  *jsonrpc.Jsonrpc
	ch   chan struct{}
	data []byte
	sync bool
}

type Fgg struct {
	pool *pool.Pool
	Conn []transport.Conn

	rpc map[string]func([]byte) (interface{}, error)

	mutex   sync.Mutex
	pending map[jsonrpc.ID]*call

	pendingMutex sync.Mutex
	pendingCount int

	finish bool

	OnSendFile func(*pool.Meta)
	OnRecvFile func(*pool.Meta)

	// Callbacks
	OnPreTran  func(*pool.Meta)
	OnPostTran func(*pool.Hash)
}

func NewFgg() *Fgg {
	fgg := &Fgg{
		pool: pool.New(),

		pendingCount: 0,

		pending:    map[jsonrpc.ID]*call{},
		OnSendFile: func(meta *pool.Meta) {},
		OnRecvFile: func(meta *pool.Meta) {},
		OnPreTran:  func(meta *pool.Meta) {},
		OnPostTran: func(meta *pool.Hash) {},
	}

	fgg.rpc = map[string]func([]byte) (interface{}, error){
		methodMeta: fgg.serverMeta,
		methodHash: fgg.serverHash,
	}

	return fgg
}

func (t *Fgg) AddConn(conn transport.Conn) {
	conn.SetOnRecv(t.recv)
	t.mutex.Lock()
	t.Conn = append(t.Conn, conn)
	t.mutex.Unlock()
}

func (t *Fgg) DelConn(conn transport.Conn) {
	remove := func(slice []transport.Conn, s int) []transport.Conn {
		return append(slice[:s], slice[s+1:]...)
	}
	t.mutex.Lock()
	defer t.mutex.Unlock()
	for i, c := range t.Conn {
		if c == conn {
			t.Conn = remove(t.Conn, i)
			return
		}
	}
}

func (t *Fgg) SetSend(file string) error {
	if err := t.pool.SetSend(file); err != nil {
		return err
	}
	meta, err := t.pool.SendMeta()
	if err != nil {
		return err
	}

	t.OnSendFile(meta)
	t.onPreTran(meta)

	data, err := json.Marshal(meta)
	if err != nil {
		log.Error(err)
	}

	t.notify(methodMeta, data)
	return err
}

func (t *Fgg) SetRecv(file string) error {
	if err := t.pool.SetRecv(file); err != nil {
		return err
	}
	t.pool.OnFinish = func() {
		t.finish = true
	}
	return nil
}

func (t *Fgg) Run(ctx context.Context) error {
	if _, err := t.rpc[methodWebrtcUp](nil); err != nil {
		log.Error(err)
	}
	ticker := time.NewTicker(loopWait)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			t.pendingMutex.Lock()
			pendingCount := t.pendingCount
			t.pendingMutex.Unlock()

			if maxPendingCount > pendingCount {
				t.getData()
			}

			if t.finish {
				if err := t.clientHash(); err != nil {
					log.Error(err)
				}
				log.Warn("run finish")
			}
		}
	}
}

func (t *Fgg) getData() {
	t.pendingMutex.Lock()
	t.pendingCount++
	t.pendingMutex.Unlock()

	c := t.pool.Next()

	if c == nil {
		time.Sleep(time.Second)
		return
	}

	data, err := json.Marshal(c)
	if err != nil {
		log.Error(err)
	}

	t.asyncCall(methodData, data)
}

func (t *Fgg) onPreTran(meta *pool.Meta) {
	t.OnPreTran(meta)
}

func (t *Fgg) onPostTran(meta *pool.Hash) {
	t.OnPostTran(meta)
}

func (f *Fgg) SetOnProgress(fn func(c int64)) {
	f.pool.OnProgress = fn
}
