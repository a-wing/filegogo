package libfgg

import (
	"encoding/json"

	"filegogo/libfgg/pool"
)

func (t *Fgg) clientHash() error {
	res, _, err := t.call(methodHash, nil)
	if err != nil {
		return err
	}

	hash := &pool.Hash{}
	if err := json.Unmarshal(res, hash); err != nil {
		return err
	}

	err = t.pool.RecvHash(hash)
	t.onPostTran(hash)
	return err
}

func (t *Fgg) serverHash(params []byte) (interface{}, error) {
	meta, err := t.pool.SendHash()
	if err != nil {
		return meta, err
	}

	t.onPostTran(meta)
	return meta, nil
}
