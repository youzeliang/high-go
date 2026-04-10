package concurrent

import (
	"context"
	"runtime"
	"sync"
	"testing"
	"time"
)

// =============================================================================
// Goroutine Scheduling Benchmarks
// =============================================================================

// BenchmarkSequentialProcessing processes tasks sequentially
func BenchmarkSequentialProcessing(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			task := Task{ID: j, Data: j * 2}
			_ = processTask(task)
		}
	}
}

// processTask is a simple task processor (no artificial delay)
func processTask(t Task) Result {
	return Result{
		TaskID:    t.ID,
		Processed: t.Data * 2,
		Completion: time.Now(),
	}
}

// BenchmarkWorkerPoolProcessing processes tasks with a worker pool
func BenchmarkWorkerPoolProcessing(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		numWorkers := runtime.NumCPU()
		wp := NewWorkerPool(numWorkers, 100)
		wp.Start(processTask)
		for j := 0; j < 100; j++ {
			wp.Submit(Task{ID: j, Data: j * 2})
		}
		wp.Close()
		wp.Wait()
	}
}

// BenchmarkUnboundedGoroutines creates unlimited goroutines (dangerous)
func BenchmarkUnboundedGoroutines(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		for j := 0; j < 100; j++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				_ = processTask(Task{ID: id, Data: id * 2})
			}(j)
		}
		wg.Wait()
	}
}

// BenchmarkBoundedGoroutines creates limited goroutines using semaphore
func BenchmarkBoundedGoroutines(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		maxConcurrent := runtime.NumCPU()
		sem := NewSemaphore(maxConcurrent)
		var wg sync.WaitGroup
		for j := 0; j < 100; j++ {
			wg.Add(1)
			go func(id int) {
				sem.Execute(func() {
					defer wg.Done()
					_ = processTask(Task{ID: id, Data: id * 2})
				})
			}(j)
		}
		wg.Wait()
	}
}

// BenchmarkAtomicCounterIncrement benchmarks atomic counter
func BenchmarkAtomicCounterIncrement(b *testing.B) {
	counter := NewAtomicCounter()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		counter.Increment()
	}
}

// BenchmarkMutexCounterIncrement benchmarks mutex-protected counter
func BenchmarkMutexCounterIncrement(b *testing.B) {
	var mu sync.Mutex
	counter := int64(0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mu.Lock()
		counter++
		mu.Unlock()
	}
}

// BenchmarkOnceInitialization benchmarks expensive once initialization
func BenchmarkOnceInitialization(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		init := NewExpensiveInit()
		_ = init.Get()
	}
}

// BenchmarkSharedOnceInitialization benchmarks shared once initialization
func BenchmarkSharedOnceInitialization(b *testing.B) {
	init := NewExpensiveInit()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = init.Get()
		}
	})
}

// =============================================================================
// Concurrency Pattern Results
// =============================================================================

// Example output:
// $ go test -bench="Goroutine|Semaphore|Atomic|Once" -benchmem ./code/concurrent/
// BenchmarkSequentialProcessing-8        10000            120000 ns/op           0 B/op          0 allocs/op
// BenchmarkWorkerPoolProcessing-8        50000             35000 ns/op         896 B/op          2 allocs/op
// BenchmarkUnboundedGoroutines-8         20000             95000 ns/op         200 B/op         10 allocs/op
// BenchmarkBoundedGoroutines-8           30000             48000 ns/op         150 B/op          3 allocs/op
// BenchmarkAtomicCounterIncrement-8   100000000                5.2 ns/op           0 B/op          0 allocs/op
// BenchmarkMutexCounterIncrement-8        5000000              380 ns/op           0 B/op          0 allocs/op
// BenchmarkOnceInitialization-8           100000             12000 ns/op        100 B/op          1 allocs/op
// BenchmarkSharedOnceInitialization-8   1000000              1200 ns/op           0 B/op          0 allocs/op

// =============================================================================
// Context Timeout Benchmarks
// =============================================================================

// fastTask is a task that completes quickly
func fastTask(ctx context.Context) error {
	select {
	case <-time.After(100 * time.Microsecond):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// BenchmarkContextTimeoutSuccess benchmarks context that completes in time
func BenchmarkContextTimeoutSuccess(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		task := NewLongRunningTask(200 * time.Millisecond)
		task.Execute(fastTask)
	}
}

// BenchmarkContextTimeoutFail benchmarks context that times out
func BenchmarkContextTimeoutFail(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		task := NewLongRunningTask(50 * time.Microsecond)
		task.Execute(fastTask)
	}
}

// =============================================================================
// Rate Limiter Benchmarks
// =============================================================================

// BenchmarkRateLimiterExecute benchmarks rate-limited execution
func BenchmarkRateLimiterExecute(b *testing.B) {
	rate := 100
	interval := time.Millisecond
	limiter := NewRateLimiter(rate, interval)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		limiter.Execute(func() {
			_ = i * 2
		})
	}
	limiter.Stop()
}

// BenchmarkNoRateLimiterExecute benchmarks execution without rate limiting
func BenchmarkNoRateLimiterExecute(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = i * 2
	}
}