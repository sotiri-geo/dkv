// Package store contains read/write logic for k-v store
package store

import (
	"errors"
	"sync"
)

var ErrKeyNotFound = errors.New("key not found error")

type KVStore struct {
	mu   sync.RWMutex
	Data map[string]string
}

func NewKVStore() *KVStore {
	return &KVStore{
		Data: make(map[string]string),
	}
}

func (s *KVStore) Put(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Data[key] = value
}

func (s *KVStore) Get(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, exists := s.Data[key]
	if !exists {
		return "", ErrKeyNotFound
	}

	return value, nil
}

func (s *KVStore) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.Data[key]; !exists {
		return ErrKeyNotFound
	}
	delete(s.Data, key)
	return nil
}
