package syncprim

import (
	"sync"
	"sync/atomic"
)

// =============================================================================
// Mutex vs RWMutex - When to use which
// =============================================================================

// CounterMutex uses a standard mutex for a counter
type CounterMutex struct {
	mu    sync.Mutex
	count int64
}

// NewCounterMutex creates a new mutex counter
func NewCounterMutex() *CounterMutex {
	return &CounterMutex{}
}

// IncMutex increments the counter (write operation)
func (c *CounterMutex) IncMutex() {
	c.mu.Lock()
	c.count++
	c.mu.Unlock()
}

// GetMutex retrieves the count (read operation)
func (c *CounterMutex) GetMutex() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

// CounterRWMutex uses a read-write mutex for a counter
type CounterRWMutex struct {
	mu    sync.RWMutex
	count int64
}

// NewCounterRWMutex creates a new RWMutex counter
func NewCounterRWMutex() *CounterRWMutex {
	return &CounterRWMutex{}
}

// IncRWMutex increments the counter (write operation)
func (c *CounterRWMutex) IncRWMutex() {
	c.mu.Lock()
	c.count++
	c.mu.Unlock()
}

// GetRWMutex retrieves the count using a read lock
func (c *CounterRWMutex) GetRWMutex() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.count
}

// CounterAtomic uses atomic operations for a counter
type CounterAtomic struct {
	count int64
}

// NewCounterAtomic creates a new atomic counter
func NewCounterAtomic() *CounterAtomic {
	return &CounterAtomic{}
}

// IncAtomic atomically increments the counter
func (c *CounterAtomic) IncAtomic() {
	atomic.AddInt64(&c.count, 1)
}

// GetAtomic atomically loads the counter
func (c *CounterAtomic) GetAtomic() int64 {
	return atomic.LoadInt64(&c.count)
}

// =============================================================================
// RWMutex - Read-heavy vs Write-heavy workloads
// =============================================================================

// ReadHeavyStore uses RWMutex for read-heavy workload
type ReadHeavyStore struct {
	mu    sync.RWMutex
	data  map[string]int64
}

// NewReadHeavyStore creates a store optimized for read-heavy workloads
func NewReadHeavyStore() *ReadHeavyStore {
	return &ReadHeavyStore{
		data: make(map[string]int64),
	}
}

// Write writes a value (exclusive lock)
func (s *ReadHeavyStore) Write(key string, value int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

// Read reads a value (shared lock)
func (s *ReadHeavyStore) Read(key string) (int64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

// WriteHeavyStore uses regular mutex for write-heavy workload
type WriteHeavyStore struct {
	mu   sync.Mutex
	data map[string]int64
}

// NewWriteHeavyStore creates a store optimized for write-heavy workloads
func NewWriteHeavyStore() *WriteHeavyStore {
	return &WriteHeavyStore{
		data: make(map[string]int64),
	}
}

// Write writes a value
func (s *WriteHeavyStore) Write(key string, value int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

// Read reads a value
func (s *WriteHeavyStore) Read(key string) (int64, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.data[key]
	return val, ok
}

// =============================================================================
// sync.Once - Single execution
// =============================================================================

// ExpensiveOperation simulates expensive initialization
type ExpensiveOperation struct {
	once     sync.Once
	result   string
	initTime int64
}

// NewExpensiveOperation creates a new expensive operation
func NewExpensiveOperation() *ExpensiveOperation {
	return &ExpensiveOperation{}
}

// GetResult gets the result, initializing only once
func (e *ExpensiveOperation) GetResult() string {
	e.once.Do(func() {
		e.result = "initialized"
		e.initTime = 1
	})
	return e.result
}

// =============================================================================
// sync.Cond - Condition variable
// =============================================================================

// SignalDemo demonstrates condition variable signaling
type SignalDemo struct {
	mu      sync.Mutex
	cond    *sync.Cond
	ready   bool
}

// NewSignalDemo creates a new signal demo
func NewSignalDemo() *SignalDemo {
	d := &SignalDemo{}
	d.cond = sync.NewCond(&d.mu)
	return d
}

// WaitForReady waits until ready is true
func (d *SignalDemo) WaitForReady() {
	d.mu.Lock()
	defer d.mu.Unlock()
	for !d.ready {
		d.cond.Wait()
	}
}

// SetReady signals that ready is true
func (d *SignalDemo) SetReady() {
	d.mu.Lock()
	d.ready = true
	d.mu.Unlock()
	d.cond.Signal()
}

// Broadcast signals all waiting goroutines
func (d *SignalDemo) Broadcast() {
	d.mu.Lock()
	d.ready = true
	d.mu.Unlock()
	d.cond.Broadcast()
}

// =============================================================================
// Atomic operations - Various atomic types
// =============================================================================

// AtomicInt64 is a wrapper for atomic int64 operations
type AtomicInt64 struct {
	value int64
}

// NewAtomicInt64 creates a new atomic int64
func NewAtomicInt64(val int64) *AtomicInt64 {
	return &AtomicInt64{value: val}
}

// Add adds a value atomically
func (a *AtomicInt64) Add(delta int64) {
	atomic.AddInt64(&a.value, delta)
}

// Load loads the value atomically
func (a *AtomicInt64) Load() int64 {
	return atomic.LoadInt64(&a.value)
}

// Store stores a value atomically
func (a *AtomicInt64) Store(val int64) {
	atomic.StoreInt64(&a.value, val)
}

// Swap swaps the value atomically
func (a *AtomicInt64) Swap(val int64) int64 {
	return atomic.SwapInt64(&a.value, val)
}

// CompareAndSwap compares and swaps the value
func (a *AtomicInt64) CompareAndSwap(old, new int64) bool {
	return atomic.CompareAndSwapInt64(&a.value, old, new)
}
