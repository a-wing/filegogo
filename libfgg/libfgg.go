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
	"github.com/sirupsen/logrus"
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
	Conn transport.Conn

	rpc map[string]func([]byte) (interface{}, error)

	mutex   sync.Mutex
	pending map[jsonrpc.ID]*call

	pendingCount int

	finish bool

	// Callbacks
	OnPreTran  func(*pool.Meta)
	OnPostTran func(*pool.Hash)
}

func NewFgg() *Fgg {
	fgg := &Fgg{
		pool: pool.New(),

		pendingCount: 0,

		pending:    map[jsonrpc.ID]*call{},
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
	t.Conn = conn
	t.Conn.SetOnRecv(t.recv)
}

func (t *Fgg) SetSend(file string) error {
	return t.pool.SetSend(file)
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
	ticker := time.NewTicker(loopWait)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if maxPendingCount > t.pendingCount {
				t.getData()
			}

			if t.finish {
				if err := t.clientHash(); err != nil {
					logrus.Error(err)
				}
				logrus.Warn("run finish")
			}
		}
	}
}

func (t *Fgg) getData() {
	t.pendingCount++

	c := t.pool.Next()

	if c == nil {
		time.Sleep(time.Second)
		return
	}

	data, err := json.Marshal(c)
	if err != nil {
		logrus.Error(err)
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
