// Package store is used to separate the operations of the store
package store

// We will implement the command design pattern to separate the
// concern of the operation from the Store struct

// Command interface for stateful operations
type Command interface {
	Apply(s *KVStore)
}

// Query interface for stateless operations
type Query interface {
	Execute(s *KVStore) (string, error)
}

// PutCommand operation to add/update values in store
type PutCommand struct {
	key   string // Keep as immutable
	value string
}

// DeleteCommand operation to remove k-v in store
type DeleteCommand struct {
	key string
}

func NewPutCommand(key, value string) *PutCommand {
	return &PutCommand{key, value}
}

func (c *PutCommand) Apply(s *KVStore) {
	s.data[c.key] = c.value
}

type GetQuery struct {
	key string
}

func NewGetQuery(key string) *GetQuery {
	return &GetQuery{key}
}

func (q *GetQuery) Execute(s *KVStore) (string, error) {
	value, exists := s.data[q.key]
	if !exists {
		return "", ErrKeyNotFound
	}
	return value, nil
}

func NewDeleteCommand(key string) *DeleteCommand {
	return &DeleteCommand{key}
}

func (c *DeleteCommand) Apply(s *KVStore) {
	delete(s.data, c.key)
}
