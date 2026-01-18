// Package store contains read/write logic for k-v store
package store

import (
	"errors"
	"fmt"
	"sync"
)

var ErrKeyNotFound = errors.New("key not found error")

type KVStore struct {
	mu   sync.RWMutex
	data map[string]string
	log  []Command
}

func NewKVStore() *KVStore {
	return &KVStore{
		data: make(map[string]string),
	}
}

func (s *KVStore) Apply(cmd Command) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Apply mutable command
	cmd.Apply(s)

	// Append to changelog
	s.log = append(s.log, cmd)
}

func (s *KVStore) Execute(query Query) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, err := query.Execute(s)
	if err != nil {
		return "", fmt.Errorf("executing query: %w", err)
	}

	return value, err
}

func (s *KVStore) GetLog() []Command {
	return s.log
}

func (s *KVStore) Replay() {
	// Reset store
	s.data = make(map[string]string)

	// Replay log commands
	for _, cmd := range s.log {
		s.Apply(cmd)
	}
}
