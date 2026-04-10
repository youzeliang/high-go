package profile

import (
	"testing"
)

// BenchmarkComputeIntensive benchmarks CPU-intensive computation
func BenchmarkComputeIntensive(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ComputeIntensive(20)
	}
}

// BenchmarkMemoryAllocation benchmarks memory allocation patterns
func BenchmarkMemoryAllocation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MemoryAllocation()
	}
}

// BenchmarkPreallocatedMemory benchmarks pre-allocated memory
func BenchmarkPreallocatedMemory(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = PreallocatedMemory()
	}
}

// BenchmarkMultipleExpensiveFunctions benchmarks calling multiple expensive functions
func BenchmarkMultipleExpensiveFunctions(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MultipleExpensiveFunctions()
	}
}

// BenchmarkMemoryIntensive benchmarks memory-intensive operation
func BenchmarkMemoryIntensive(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MemoryIntensive(1024 * 1024) // 1MB
	}
}

// BenchmarkGetMemStats benchmarks reading memory stats
func BenchmarkGetMemStats(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GetMemStats()
	}
}

// BenchmarkForceGC benchmarks forced garbage collection
func BenchmarkForceGC(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ForceGC()
	}
}

// Example profiling commands:
// $ go test -bench="ComputeIntensive" -cpuprofile=cpu.prof ./code/profile/
// $ go tool pprof cpu.prof
// (pprof) top
// (pprof) web
//
// $ go test -bench="MemoryAllocation" -memprofile=mem.prof ./code/profile/
// $ go tool pprof mem.prof
//
// $ go test -bench="." -trace=trace.out ./code/profile/
// $ go tool trace trace.out
