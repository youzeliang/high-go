package concurrent

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// =============================================================================
// Goroutine Efficient Scheduling - Worker Pool Pattern
// =============================================================================

// Task represents a unit of work
type Task struct {
	ID   int
	Data int
}

// Result represents the result of a task
type Result struct {
	TaskID     int
	Processed  int
	Completion time.Time
}

// WorkerPool represents a pool of workers for concurrent task processing
type WorkerPool struct {
	tasks   chan Task
	results chan Result
	workers int
	wg      sync.WaitGroup
}

// NewWorkerPool creates a new worker pool with specified number of workers
func NewWorkerPool(workers, taskBuffer int) *WorkerPool {
	return &WorkerPool{
		tasks:   make(chan Task, taskBuffer),
		results: make(chan Result, taskBuffer),
		workers: workers,
	}
}

// Start starts the worker pool
func (wp *WorkerPool) Start(handler func(Task) Result) {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go func(workerID int) {
			defer wp.wg.Done()
			for task := range wp.tasks {
				result := handler(task)
				wp.results <- result
			}
		}(i)
	}
}

// Submit submits a task to the pool
func (wp *WorkerPool) Submit(task Task) {
	wp.tasks <- task
}

// Close closes the task channel
func (wp *WorkerPool) Close() {
	close(wp.tasks)
}

// Wait waits for all workers to complete
func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
	close(wp.results)
}

// Results returns the results channel
func (wp *WorkerPool) Results() <-chan Result {
	return wp.results
}

// =============================================================================
// Bounded Goroutines - Limit concurrent goroutines to prevent resource exhaustion
// =============================================================================

// Semaphore implements a semaphore pattern to limit concurrent operations
type Semaphore struct {
	tokens chan struct{}
}

// NewSemaphore creates a new semaphore with the specified max concurrent operations
func NewSemaphore(maxConcurrent int) *Semaphore {
	s := &Semaphore{
		tokens: make(chan struct{}, maxConcurrent),
	}
	// Initialize with tokens
	for i := 0; i < maxConcurrent; i++ {
		s.tokens <- struct{}{}
	}
	return s
}

// Acquire acquires a token (blocks if unavailable)
func (s *Semaphore) Acquire() {
	<-s.tokens
}

// Release releases a token back to the semaphore
func (s *Semaphore) Release() {
	s.tokens <- struct{}{}
}

// Execute runs the function with semaphore protection
func (s *Semaphore) Execute(fn func()) {
	s.Acquire()
	defer s.Release()
	fn()
}

// =============================================================================
// Goroutine Pool with GOMAXPROCS - Proper CPU utilization
// =============================================================================

// ConfiguredWorkerPool is a worker pool that respects GOMAXPROCS
type ConfiguredWorkerPool struct {
	pool   sync.Pool
	ctx    context.Context
	cancel context.CancelFunc
}

// NewConfiguredWorkerPool creates a worker pool optimized for CPU count
func NewConfiguredWorkerPool() *ConfiguredWorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	numCPU := runtime.NumCPU()
	pool := sync.Pool{
		New: func() interface{} {
			return &struct {
				tasks   chan func() interface{}
				results chan interface{}
				done    chan struct{}
			}{
				tasks:   make(chan func() interface{}, numCPU),
				results: make(chan interface{}, numCPU),
				done:    make(chan struct{}),
			}
		},
	}
	return &ConfiguredWorkerPool{
		pool:   pool,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Process submits work and returns a channel for the result
func (cwp *ConfiguredWorkerPool) Process(task func() interface{}) <-chan interface{} {
	worker := cwp.pool.Get().(*struct {
		tasks   chan func() interface{}
		results chan interface{}
		done    chan struct{}
	})
	resultCh := make(chan interface{})
	go func() {
		result := task()
		resultCh <- result
		cwp.pool.Put(worker)
	}()
	return resultCh
}

// Shutdown shuts down the worker pool
func (cwp *ConfiguredWorkerPool) Shutdown() {
	cwp.cancel()
}

// =============================================================================
// Rate Limiter - Token bucket algorithm for controlling goroutine spawn rate
// =============================================================================

// RateLimiter implements a token bucket rate limiter
type RateLimiter struct {
	tokens    chan struct{}
	rate      int
	interval  time.Duration
	stopCh    chan struct{}
	wg        sync.WaitGroup
}

// NewRateLimiter creates a rate limiter that allows 'rate' operations per interval
func NewRateLimiter(rate int, interval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		tokens:   make(chan struct{}, rate),
		rate:     rate,
		interval: interval,
		stopCh:   make(chan struct{}),
	}
	// Fill the bucket
	for i := 0; i < rate; i++ {
		rl.tokens <- struct{}{}
	}
	// Start token replenishment
	rl.wg.Add(1)
	go func() {
		defer rl.wg.Done()
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				select {
				case rl.tokens <- struct{}{}:
				default:
				}
			case <-rl.stopCh:
				return
			}
		}
	}()
	return rl
}

// Wait waits for permission to proceed (blocks until token available)
func (rl *RateLimiter) Wait() {
	<-rl.tokens
	// Don't put token back - it's consumed
}

// Stop stops the rate limiter
func (rl *RateLimiter) Stop() {
	close(rl.stopCh)
	rl.wg.Wait()
}

// Execute runs the function with rate limiting
func (rl *RateLimiter) Execute(fn func()) {
	rl.Wait()
	fn()
}

// =============================================================================
// Atomic Counter - Lock-free concurrent counter
// =============================================================================

// AtomicCounter is a lock-free concurrent counter
type AtomicCounter struct {
	count int64
}

// NewAtomicCounter creates a new atomic counter
func NewAtomicCounter() *AtomicCounter {
	return &AtomicCounter{}
}

// Increment increments the counter
func (ac *AtomicCounter) Increment() {
	atomic.AddInt64(&ac.count, 1)
}

// Decrement decrements the counter
func (ac *AtomicCounter) Decrement() {
	atomic.AddInt64(&ac.count, -1)
}

// Value returns the current value
func (ac *AtomicCounter) Value() int64 {
	return atomic.LoadInt64(&ac.count)
}

// Add adds a value to the counter
func (ac *AtomicCounter) Add(v int64) {
	atomic.AddInt64(&ac.count, v)
}

// =============================================================================
// Once Initialization - Ensure expensive initialization happens only once
// =============================================================================

// ExpensiveInit simulates expensive initialization
type ExpensiveInit struct {
	once     sync.Once
	value    string
	initTime time.Time
}

// NewExpensiveInit creates a new expensive initialization struct
func NewExpensiveInit() *ExpensiveInit {
	return &ExpensiveInit{}
}

// Get returns the initialized value (initializes only once)
func (e *ExpensiveInit) Get() string {
	e.once.Do(func() {
		time.Sleep(10 * time.Millisecond) // Simulate expensive operation
		e.value = "initialized"
		e.initTime = time.Now()
	})
	return e.value
}

// =============================================================================
// Context with Timeout - Proper goroutine cancellation
// =============================================================================

// LongRunningTask represents a long-running task that can be cancelled
type LongRunningTask struct {
	ctx    context.Context
	cancel context.CancelFunc
}

// NewLongRunningTask creates a new long-running task with timeout
func NewLongRunningTask(timeout time.Duration) *LongRunningTask {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return &LongRunningTask{
		ctx:    ctx,
		cancel: cancel,
	}
}

// Execute runs the task, returning false if cancelled or timed out
func (lrt *LongRunningTask) Execute(fn func(context.Context) error) bool {
	done := make(chan error, 1)
	go func() {
		done <- fn(lrt.ctx)
	}()
	select {
	case err := <-done:
		return err == nil
	case <-lrt.ctx.Done():
		return false
	}
}

// Cancel cancels the task
func (lrt *LongRunningTask) Cancel() {
	lrt.cancel()
}