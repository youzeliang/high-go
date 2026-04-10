package mapconv

import (
	"sync"
)

// MutexMap uses a mutex to protect map access
type MutexMap struct {
	mu   sync.Mutex
	data map[string]int
}

// NewMutexMap creates a new mutex-protected map
func NewMutexMap() *MutexMap {
	return &MutexMap{
		data: make(map[string]int),
	}
}

// Get retrieves a value from the map
func (m *MutexMap) Get(key string) (int, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	val, ok := m.data[key]
	return val, ok
}

// Set stores a value in the map
func (m *MutexMap) Set(key string, value int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

// Delete removes a key from the map
func (m *MutexMap) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
}

// Len returns the number of elements in the map
func (m *MutexMap) Len() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.data)
}

// RWMutexMap uses read-write mutex for better read performance
type RWMutexMap struct {
	mu   sync.RWMutex
	data map[string]int
}

// NewRWMutexMap creates a new RWMutex-protected map
func NewRWMutexMap() *RWMutexMap {
	return &RWMutexMap{
		data: make(map[string]int),
	}
}

// Get retrieves a value using read lock
func (m *RWMutexMap) Get(key string) (int, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.data[key]
	return val, ok
}

// Set stores a value using write lock
func (m *RWMutexMap) Set(key string, value int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

// Delete removes a key using write lock
func (m *RWMutexMap) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
}

// Len returns the number of elements using read lock
func (m *RWMutexMap) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.data)
}

// SyncMap uses Go's sync.Map for concurrent access
type SyncMap struct {
	data sync.Map
}

// NewSyncMap creates a new sync.Map wrapper
func NewSyncMap() *SyncMap {
	return &SyncMap{}
}

// Get retrieves a value from sync.Map
func (m *SyncMap) Get(key string) (int, bool) {
	val, ok := m.data.Load(key)
	if !ok {
		return 0, false
	}
	return val.(int), true
}

// Set stores a value in sync.Map
func (m *SyncMap) Set(key string, value int) {
	m.data.Store(key, value)
}

// Delete removes a key from sync.Map
func (m *SyncMap) Delete(key string) {
	m.data.Delete(key)
}

// Len returns the number of elements (approximate)
func (m *SyncMap) Len() int {
	count := 0
	m.data.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return count
}

// BulkSyncMap loads multiple key-value pairs efficiently
func (m *SyncMap) BulkSyncMap(pairs map[string]int) {
	for k, v := range pairs {
		m.data.Store(k, v)
	}
}

// BatchMutexMap demonstrates batch operations with mutex
type BatchMutexMap struct {
	mu   sync.Mutex
	data map[string]int
}

// NewBatchMutexMap creates a new batch mutex map
func NewBatchMutexMap() *BatchMutexMap {
	return &BatchMutexMap{
		data: make(map[string]int),
	}
}

// SetBatch stores multiple values in one lock acquisition
func (m *BatchMutexMap) SetBatch(pairs map[string]int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range pairs {
		m.data[k] = v
	}
}

// GetBatch retrieves multiple values
func (m *BatchMutexMap) GetBatch(keys []string) map[string]int {
	result := make(map[string]int, len(keys))
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, k := range keys {
		if v, ok := m.data[k]; ok {
			result[k] = v
		}
	}
	return result
}
