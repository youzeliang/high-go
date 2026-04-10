package profile

import (
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"time"
)

// =============================================================================
// CPU Profiling Examples
// =============================================================================

// ExpensiveFunction simulates expensive CPU work
func ExpensiveFunction() int {
	sum := 0
	for i := 0; i < 1000000; i++ {
		sum += i
	}
	return sum
}

// MultipleExpensiveFunctions calls expensive functions multiple times
func MultipleExpensiveFunctions() int {
	sum := 0
	for i := 0; i < 100; i++ {
		sum += ExpensiveFunction()
	}
	return sum
}

// MemoryAllocation demonstrates memory allocations
func MemoryAllocation() []int {
	// This creates multiple allocations
	result := make([]int, 0)
	for i := 0; i < 1000; i++ {
		result = append(result, i)
	}
	return result
}

// PreallocatedMemory demonstrates efficient memory usage
func PreallocatedMemory() []int {
	// This creates only one allocation
	result := make([]int, 0, 1000)
	for i := 0; i < 1000; i++ {
		result = append(result, i)
	}
	return result
}

// =============================================================================
// pprof Setup Functions (for educational purposes)
// =============================================================================

// StartCPUProfile starts CPU profiling
func StartCPUProfile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	pprof.StartCPUProfile(f)
	// In real usage, call pprof.StopCPUProfile() after work
	return nil
}

// StopCPUProfile stops CPU profiling
func StopCPUProfile() {
	pprof.StopCPUProfile()
}

// WriteMemoryProfile writes a memory profile
func WriteMemoryProfile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	pprof.WriteHeapProfile(f)
	f.Close()
	return nil
}

// WriteGoroutineProfile writes goroutine profile
func WriteGoroutineProfile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	pprof.Lookup("goroutine").WriteTo(f, 0)
	f.Close()
	return nil
}

// =============================================================================
// trace Examples
// =============================================================================

// StartTrace starts execution tracing
func StartTrace(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	if err := trace.Start(f); err != nil {
		return err
	}
	return nil
}

// StopTrace stops execution tracing
func StopTrace() {
	trace.Stop()
}

// =============================================================================
// Benchmark-style functions for profiling
// =============================================================================

// ComputeIntensive is a CPU-intensive computation
func ComputeIntensive(n int) int {
	if n <= 1 {
		return n
	}
	return ComputeIntensive(n-1) + ComputeIntensive(n-2)
}

// MemoryIntensive demonstrates memory pressure
func MemoryIntensive(size int) []byte {
	data := make([]byte, size)
	for i := 0; i < size; i++ {
		data[i] = byte(i % 256)
	}
	return data
}

// ConcurrentWorkload demonstrates concurrent workload
func ConcurrentWorkload(done chan<- int) {
	sum := 0
	for i := 0; i < 1000000; i++ {
		sum += i
	}
	done <- sum
}

// GCSettings demonstrates GC pressure
func GCSettings() {
	// Force GC to see its impact
	runtime.GC()
}

// =============================================================================
// Performance measurement helpers
// =============================================================================

// MeasureTime measures execution time
func MeasureTime(fn func()) time.Duration {
	start := time.Now()
	fn()
	return time.Since(start)
}

// GetMemStats returns current memory statistics
func GetMemStats() runtime.MemStats {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	return stats
}

// ForceGC forces garbage collection
func ForceGC() {
	runtime.GC()
}
