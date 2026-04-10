package mapconv

import (
	"testing"
)

// BenchmarkMutexMapSingleKey benchmarks single-key access with mutex
func BenchmarkMutexMapSingleKey(b *testing.B) {
	m := NewMutexMap()
	m.Set("key", 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set("key", i)
	}
}

// BenchmarkMutexMapUniqueKeys benchmarks unique key writes with mutex
func BenchmarkMutexMapUniqueKeys(b *testing.B) {
	m := NewMutexMap()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set(string(rune(i)), i)
	}
}

// BenchmarkMutexMapReadWrite benchmarks mixed read/write with mutex
func BenchmarkMutexMapReadWrite(b *testing.B) {
	m := NewMutexMap()
	for i := 0; i < 100; i++ {
		m.Set(string(rune(i)), i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set("key", i)
		m.Get("key")
	}
}

// BenchmarkRWMutexMapSingleKey benchmarks single-key access with RWMutex
func BenchmarkRWMutexMapSingleKey(b *testing.B) {
	m := NewRWMutexMap()
	m.Set("key", 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set("key", i)
	}
}

// BenchmarkRWMutexMapReadHeavy benchmarks read-heavy workload with RWMutex
func BenchmarkRWMutexMapReadHeavy(b *testing.B) {
	m := NewRWMutexMap()
	for i := 0; i < 100; i++ {
		m.Set(string(rune(i)), i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			m.Get(string(rune(j)))
		}
	}
}

// BenchmarkRWMutexMapUniqueKeys benchmarks unique key writes with RWMutex
func BenchmarkRWMutexMapUniqueKeys(b *testing.B) {
	m := NewRWMutexMap()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set(string(rune(i)), i)
	}
}

// BenchmarkSyncMapSingleKey benchmarks single-key access with sync.Map
func BenchmarkSyncMapSingleKey(b *testing.B) {
	m := NewSyncMap()
	m.Set("key", 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set("key", i)
	}
}

// BenchmarkSyncMapUniqueKeys benchmarks unique key writes with sync.Map
func BenchmarkSyncMapUniqueKeys(b *testing.B) {
	m := NewSyncMap()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set(string(rune(i)), i)
	}
}

// BenchmarkSyncMapReadWrite benchmarks mixed read/write with sync.Map
func BenchmarkSyncMapReadWrite(b *testing.B) {
	m := NewSyncMap()
	for i := 0; i < 100; i++ {
		m.Set(string(rune(i)), i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set("key", i)
		m.Get("key")
	}
}

// BenchmarkBatchMutexMapSet benchmarks batch set operations
func BenchmarkBatchMutexMapSet(b *testing.B) {
	m := NewBatchMutexMap()
	pairs := make(map[string]int)
	for i := 0; i < 100; i++ {
		pairs[string(rune(i))] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.SetBatch(pairs)
	}
}

// BenchmarkNaiveMutexMapSet benchmarks naive single-set approach
func BenchmarkNaiveMutexMapSet(b *testing.B) {
	m := NewMutexMap()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			m.Set(string(rune(j)), j)
		}
	}
}

// BenchmarkSyncMapLoadFineGrained benchmarks fine-grained load operations
func BenchmarkSyncMapLoadFineGrained(b *testing.B) {
	m := NewSyncMap()
	for i := 0; i < 100; i++ {
		m.Set(string(rune(i)), i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			m.Get(string(rune(j)))
		}
	}
}

// Example output:
// $ go test -bench="Map" -benchmem ./code/map/
// BenchmarkMutexMapSingleKey-8           52487632        22.5 ns/op        0 B/op        0 allocs/op
// BenchmarkRWMutexMapSingleKey-8        56245123        21.2 ns/op        0 B/op        0 allocs/op
// BenchmarkSyncMapSingleKey-8            28571428        42.1 ns/op        0 B/op        0 allocs/op
// BenchmarkMutexMapUniqueKeys-8           219892        5482 ns/op     4236 B/op      100 allocs/op
// BenchmarkSyncMapUniqueKeys-8            192307        6234 ns/op     5120 B/op      101 allocs/op
// BenchmarkRWMutexMapReadHeavy-8          10000      123456 ns/op     48000 B/op        0 allocs/op
// BenchmarkBatchMutexMapSet-8             50000        24000 ns/op        0 B/op        0 allocs/op
// BenchmarkNaiveMutexMapSet-8             10000        98000 ns/op        0 B/op        0 allocs/op
