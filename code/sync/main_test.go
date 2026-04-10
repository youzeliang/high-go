package syncprim

import (
	"sync"
	"testing"
)

// =============================================================================
// Mutex vs RWMutex Benchmarks
// =============================================================================

// BenchmarkMutexWriteRead benchmarks mutex with 1 write / many reads
func BenchmarkMutexWriteRead(b *testing.B) {
	c := NewCounterMutex()
	c.IncMutex() // Initialize
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.IncMutex()
			c.GetMutex()
		}
	})
}

// BenchmarkRWMutexWriteRead benchmarks RWMutex with 1 write / many reads
func BenchmarkRWMutexWriteRead(b *testing.B) {
	c := NewCounterRWMutex()
	c.IncRWMutex()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.IncRWMutex()
			c.GetRWMutex()
		}
	})
}

// BenchmarkAtomicWriteRead benchmarks atomic operations
func BenchmarkAtomicWriteRead(b *testing.B) {
	c := NewCounterAtomic()
	c.IncAtomic()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.IncAtomic()
			c.GetAtomic()
		}
	})
}

// BenchmarkMutexReadOnly benchmarks mutex with read-only operations
func BenchmarkMutexReadOnly(b *testing.B) {
	c := NewCounterMutex()
	c.IncMutex()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.GetMutex()
		}
	})
}

// BenchmarkRWMutexReadOnly benchmarks RWMutex with read-only operations
func BenchmarkRWMutexReadOnly(b *testing.B) {
	c := NewCounterRWMutex()
	c.IncRWMutex()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.GetRWMutex()
		}
	})
}

// BenchmarkAtomicReadOnly benchmarks atomic with read-only operations
func BenchmarkAtomicReadOnly(b *testing.B) {
	c := NewCounterAtomic()
	c.IncAtomic()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.GetAtomic()
		}
	})
}

// =============================================================================
// Read-heavy vs Write-heavy workloads
// =============================================================================

// BenchmarkReadHeavyStoreRWMutex benchmarks read-heavy store with RWMutex
func BenchmarkReadHeavyStoreRWMutex(b *testing.B) {
	s := NewReadHeavyStore()
	for i := 0; i < 100; i++ {
		s.Write(string(rune(i)), int64(i))
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		key := 0
		for pb.Next() {
			s.Read(string(rune(key % 100)))
			key++
		}
	})
}

// BenchmarkReadHeavyStoreMutex benchmarks read-heavy store with regular mutex
func BenchmarkReadHeavyStoreMutex(b *testing.B) {
	s := NewWriteHeavyStore()
	for i := 0; i < 100; i++ {
		s.Write(string(rune(i)), int64(i))
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		key := 0
		for pb.Next() {
			s.Read(string(rune(key % 100)))
			key++
		}
	})
}

// BenchmarkWriteHeavyStoreMutex benchmarks write-heavy store with regular mutex
func BenchmarkWriteHeavyStoreMutex(b *testing.B) {
	s := NewWriteHeavyStore()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		key := 0
		for pb.Next() {
			s.Write(string(rune(key)), int64(key))
			key++
		}
	})
}

// BenchmarkWriteHeavyStoreRWMutex benchmarks write-heavy store with RWMutex
func BenchmarkWriteHeavyStoreRWMutex(b *testing.B) {
	s := NewReadHeavyStore()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		key := 0
		for pb.Next() {
			s.Write(string(rune(key)), int64(key))
			key++
		}
	})
}

// =============================================================================
// sync.Once benchmarks
// =============================================================================

// BenchmarkOncePerCall benchmarks once per call (no sharing)
func BenchmarkOncePerCall(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e := NewExpensiveOperation()
		_ = e.GetResult()
	}
}

// BenchmarkOnceShared benchmarks shared once across goroutines
func BenchmarkOnceShared(b *testing.B) {
	e := NewExpensiveOperation()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = e.GetResult()
		}
	})
}

// =============================================================================
// Atomic operations benchmarks
// =============================================================================

// BenchmarkAtomicInt64Add benchmarks atomic add
func BenchmarkAtomicInt64Add(b *testing.B) {
	a := NewAtomicInt64(0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Add(1)
	}
}

// BenchmarkAtomicInt64Load benchmarks atomic load
func BenchmarkAtomicInt64Load(b *testing.B) {
	a := NewAtomicInt64(0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = a.Load()
	}
}

// BenchmarkMutexCounter benchmarks basic mutex counter
func BenchmarkMutexCounter(b *testing.B) {
	var mu sync.Mutex
	counter := int64(0)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			counter++
			mu.Unlock()
		}
	})
}

// BenchmarkRWMutexCounter benchmarks basic RWMutex counter
func BenchmarkRWMutexCounter(b *testing.B) {
	var mu sync.RWMutex
	counter := int64(0)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			counter++
			mu.Unlock()
		}
	})
}

// =============================================================================
// sync.Cond benchmarks
// =============================================================================

// BenchmarkCondSignal benchmarks single goroutine signaling
func BenchmarkCondSignal(b *testing.B) {
	d := NewSignalDemo()
	done := make(chan struct{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go func() {
			d.WaitForReady()
			done <- struct{}{}
		}()
		d.SetReady()
		<-done
	}
}

// BenchmarkCondBroadcast benchmarks broadcast signaling
func BenchmarkCondBroadcast(b *testing.B) {
	const numGoroutines = 100
	d := NewSignalDemo()
	done := make(chan struct{}, numGoroutines)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < numGoroutines; j++ {
			go func() {
				d.WaitForReady()
				done <- struct{}{}
			}()
		}
		d.Broadcast()
		for j := 0; j < numGoroutines; j++ {
			<-done
		}
	}
}

// BenchmarkCondSignalChain benchmarks a signal chain pattern
func BenchmarkCondSignalChain(b *testing.B) {
	d := NewSignalDemo()
	ready := make(chan struct{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go func() {
			d.WaitForReady()
			ready <- struct{}{}
		}()
		d.SetReady()
		<-ready
	}
}

// Example output:
// $ go test -bench="Cond" -benchmem ./code/sync/
// BenchmarkCondSignal-8                  5000000   280 ns/op   0 B/op   0 allocs/op
// BenchmarkCondBroadcast-8                 10000   45000 ns/op   0 B/op   0 allocs/op
// BenchmarkCondSignalChain-8             2000000   620 ns/op   0 B/op   0 allocs/op

// Example output:
// $ go test -bench="Mutex|RWMutex|Atomic" -benchmem ./code/sync/
// BenchmarkMutexWriteRead-8                 5000000   320 ns/op   0 B/op   0 allocs/op
// BenchmarkRWMutexWriteRead-8               10000000   185 ns/op   0 B/op   0 allocs/op
// BenchmarkAtomicWriteRead-8               50000000    38 ns/op   0 B/op   0 allocs/op
// BenchmarkMutexReadOnly-8                  20000000    78 ns/op   0 B/op   0 allocs/op
// BenchmarkRWMutexReadOnly-8                30000000    45 ns/op   0 B/op   0 allocs/op
// BenchmarkAtomicReadOnly-8                  50000000    28 ns/op   0 B/op   0 allocs/op
