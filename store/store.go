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
		log:  make([]Command, 0),
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
	// Add read lock
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Returning a copy to prevent external modification
	newLog := make([]Command, len(s.log))
	copy(newLog, s.log)
	return newLog
}

func (s *KVStore) Replay() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[string]string)

	// Replay log commands, do not add to the log
	for _, cmd := range s.log {
		cmd.Apply(s)
	}
}
