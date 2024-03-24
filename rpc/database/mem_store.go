package database

import (
	"sync"

	"github.com/google/uuid"
)

type memStore struct {
	Requests      map[uuid.UUID]RequestEntry
	EthSendRawTxs map[uuid.UUID][]*EthSendRawTxEntry
	mutex         sync.Mutex
}

func NewMemStore() *memStore {
	return &memStore{
		Requests:      make(map[uuid.UUID]RequestEntry),
		EthSendRawTxs: make(map[uuid.UUID][]*EthSendRawTxEntry),
		mutex:         sync.Mutex{},
	}
}

func (m *memStore) SaveRequestEntry(entry RequestEntry) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Requests[entry.ID] = entry
	return nil
}

func (m *memStore) SaveRawTxEntries(entries []*EthSendRawTxEntry) error {
	if len(entries) != 0 {
		m.mutex.Lock()
		defer m.mutex.Unlock()
		m.EthSendRawTxs[entries[0].RequestID] = entries
	}
	return nil
}
