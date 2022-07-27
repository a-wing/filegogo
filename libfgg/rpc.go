package libfgg

import (
	"encoding/json"
	"errors"
	"filegogo/libfgg/pool"

	"github.com/sb-im/jsonrpc-lite"
	log "github.com/sirupsen/logrus"
)

func (t *Fgg) send(head []byte, body []byte) error {
	log.Trace(string(head), len(body))
	t.mutex.Lock()
	if n := len(t.Conn); n > 0 {
		c := t.Conn[n-1]
		t.mutex.Unlock()
		return c.Send(head, body)
	} else {
		t.mutex.Unlock()
	}
	return errors.New("Not found conn")
}

func (t *Fgg) recv(head []byte, body []byte) {
	log.Trace(string(head), len(body))
	rpc := jsonrpc.ParseObject(head)

	var resHead []byte
	var resBody []byte

	if rpc.Type == jsonrpc.TypeRequest || rpc.Type == jsonrpc.TypeNotify {

		if fn, ok := t.rpc[rpc.Method]; ok {
			var params []byte
			if rpc.Params != nil {
				params = *rpc.Params
			} else {
				params = nil
			}
			res, err := fn(params)
			if rpc.Type == jsonrpc.TypeNotify {
				if err != nil {
					log.Error(err)
				}
				return
			}

			if err != nil {
				resHead, _ = jsonrpc.NewError(rpc.ID, 404, err.Error(), nil).ToJSON()
			} else {
				resHead, _ = jsonrpc.NewSuccess(rpc.ID, res).ToJSON()
			}
			//fmt.Printf("%+v\n", rpc)
		} else if rpc.Method == methodData {
			//fmt.Println("server getdata")
			_, resBody, _ = t.serverSendData(*rpc.Params)
			resHead, _ = jsonrpc.NewSuccess(rpc.ID, *rpc.Params).ToJSON()
		} else {
			resHead, _ = jsonrpc.NewError(rpc.ID, 404, "Not Found this Method", nil).ToJSON()
		}
		t.send(resHead, resBody)
	} else if rpc.Type == jsonrpc.TypeSuccess || rpc.Type == jsonrpc.TypeErrors {
		t.mutex.Lock()
		cc, ok := t.pending[*rpc.ID]
		t.mutex.Unlock()
		if !ok {
			return
		}

		if cc.req.Method == methodData {
			t.pendingMutex.Lock()
			t.pendingCount--
			t.pendingMutex.Unlock()

			c := &pool.DataChunk{}
			json.Unmarshal(*cc.req.Params, c)

			t.pool.RecvData(c, body)
		} else {
			if cc.sync {
				cc.res = rpc
				cc.data = body
				cc.ch <- struct{}{}
			}
		}
	} else {
		log.Errorf("Unknown rpc: %+v\n", rpc)
	}
}

func (t *Fgg) serverSendData(raw []byte) ([]byte, []byte, error) {
	c := &pool.DataChunk{}
	json.Unmarshal(raw, c)

	data, err := t.pool.SendData(c)
	return nil, data, err
}

func (t *Fgg) call(method string, params []byte) ([]byte, []byte, error) {
	wait, err := t.doCall(method, params, true)
	if err != nil {
		return nil, nil, err
	}

	<-wait.ch
	res := wait.res
	if res.Type == jsonrpc.TypeSuccess {
		return *res.Result, wait.data, nil
	} else if res.Type == jsonrpc.TypeErrors {
		return nil, wait.data, errors.New(res.Errors.Message)
	}

	return nil, wait.data, errors.New("rpc error")
}

func (t *Fgg) asyncCall(method string, params []byte) (*call, error) {
	return t.doCall(method, params, false)
}

func (t *Fgg) doCall(method string, params []byte, sync bool) (*call, error) {
	rpc := jsonrpc.NewRequest(getUniqueID(), method, json.RawMessage(params))
	data, err := rpc.ToJSON()
	if err != nil {
		return nil, err
	}

	t.mutex.Lock()

	c := &call{
		req:  rpc,
		ch:   make(chan struct{}, 1),
		sync: sync,
	}
	t.pending[*rpc.ID] = c
	t.mutex.Unlock()
	if err := t.send(data, nil); err != nil {
		t.mutex.Lock()
		delete(t.pending, *rpc.ID)
		t.mutex.Unlock()
		return c, err
	}

	return c, nil
}

func (t *Fgg) notify(method string, params []byte) error {
	data, err := jsonrpc.NewNotify(method, json.RawMessage(params)).ToJSON()
	if err != nil {
		return err
	}
	return t.send(data, nil)
}
