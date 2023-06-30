package store

import (
	"encoding/json"
	"sync"
)

type Store struct {
	mutex sync.Mutex
	store map[string][]byte
}

func NewStore() *Store {
	return &Store{
		mutex: sync.Mutex{},
		store: make(map[string][]byte),
	}
}

func (s *Store) Put(key string, value any) error {
	data, err := json.Marshal(value)
	s.mutex.Lock()
	s.store[key] = data
	s.mutex.Unlock()
	return err
}

func (s *Store) Get(key string, value any) error {
	s.mutex.Lock()
	data := s.store[key]
	s.mutex.Unlock()
	return json.Unmarshal(data, value)
}

func (s *Store) Del(key string) {
	s.mutex.Lock()
	delete(s.store, key)
	s.mutex.Unlock()
}
