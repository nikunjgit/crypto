package store

import (
	"sync"
	"github.com/nikunjgit/crypto/event"
)

type MemoryStore struct {
	memory map[string]event.Messages
	ttl    int
	mutex  *sync.Mutex
}

func NewMemoryStore(ttl int) *MemoryStore {
	return &MemoryStore{make(map[string]event.Messages), ttl, &sync.Mutex{}}
}

func (m *MemoryStore) Set(key string, message event.Messages) error {
	m.mutex.Lock()
	m.memory[key] = message
	m.mutex.Unlock()
	return nil
}

func (m *MemoryStore) Get(keys []string) (event.Messages, error) {
	vals := make(event.Messages, 0, 10)
	m.mutex.Lock()
	for _, key := range keys {
		val, ok := m.memory[key]
		if !ok {
			continue
		}
		vals = append(vals, val...)
	}

	m.mutex.Unlock()
	return vals, nil
}
