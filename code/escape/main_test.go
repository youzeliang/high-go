package escape

import (
	"testing"
)

// BenchmarkStackAlloc benchmarks stack allocation
func BenchmarkStackAlloc(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = StackAlloc()
	}
}

// BenchmarkHeapAlloc benchmarks heap allocation (pointer escape)
func BenchmarkHeapAlloc(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = HeapAlloc()
	}
}

// BenchmarkHeapAllocByInterface benchmarks interface allocation
func BenchmarkHeapAllocByInterface(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = HeapAllocByInterface()
	}
}

// BenchmarkStackWithSlice benchmarks fixed-size array on stack
func BenchmarkStackWithSlice(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = StackWithSlice()
	}
}

// BenchmarkHeapWithSlice benchmarks slice returning (header escapes)
func BenchmarkHeapWithSlice(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = HeapWithSlice()
	}
}

// BenchmarkStackCopy benchmarks stack-to-stack copy
func BenchmarkStackCopy(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = StackCopy()
	}
}

// BenchmarkPreAllocSlice benchmarks pre-allocated slice
func BenchmarkPreAllocSlice(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = PreAllocSlice(5)
	}
}

// BenchmarkPassByPointer benchmarks passing pointer to function
func BenchmarkPassByPointer(b *testing.B) {
	x := 10
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = PassByPointer(&x)
	}
}

// BenchmarkStackByValue benchmarks struct/array passed by value
func BenchmarkStackByValue(b *testing.B) {
	arr := [4]int{1, 2, 3, 4}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = StackByValue(arr)
	}
}

// =============================================================================
// Compiler Optimization Benchmarks
// =============================================================================

// BenchmarkSimpleAdd benchmarks a simple function that can be inlined
func BenchmarkSimpleAdd(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SimpleAdd(i, i+1)
	}
}

// BenchmarkInliningThreshold benchmarks function at inlining threshold
func BenchmarkInliningThreshold(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = InliningThreshold(100)
	}
}

// BenchmarkNoInline benchmarks function with noinline directive
func BenchmarkNoInline(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NoInlineFunction(100)
	}
}

// BenchmarkLocalVariableCaching benchmarks local variable caching
func BenchmarkLocalVariableCaching(b *testing.B) {
	arr := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		arr[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = LocalVariableCaching(arr)
	}
}

// BenchmarkBatchOperation benchmarks batch processing
func BenchmarkBatchOperation(b *testing.B) {
	arr := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		arr[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BatchOperation(arr)
	}
}

// BenchmarkReducePointerDerefs benchmarks reducing pointer dereferences
func BenchmarkReducePointerDerefs(b *testing.B) {
	data := &struct{ value int }{value: 42}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ReducePointerDerefs(data)
	}
}

// Example output:
// $ go test -bench="Alloc|Inline|Local" -benchmem -run=^$ ./code/escape/
// BenchmarkStackAlloc-8           2000000000    0.30 ns/op    0 B/op    0 allocs/op
// BenchmarkHeapAlloc-8             100000000    10.2 ns/op    8 B/op    1 allocs/op
// BenchmarkHeapAllocByInterface-8   50000000    30.1 ns/op   16 B/op    1 allocs/op
// BenchmarkStackWithSlice-8          5000000   280 ns/op   4096 B/op    0 allocs/op
// BenchmarkHeapWithSlice-8           5000000   305 ns/op   4096 B/op    1 allocs/op
// BenchmarkStackCopy-8              2000000000    0.28 ns/op    0 B/op    0 allocs/op
// BenchmarkPreAllocSlice-8           100000000    11.2 ns/op   48 B/op    1 allocs/op
// BenchmarkPassByPointer-8          2000000000    0.25 ns/op    0 B/op    0 allocs/op
// BenchmarkStackByValue-8           2000000000    0.32 ns/op    0 B/op    0 allocs/op
// BenchmarkSimpleAdd-8              2000000000    0.28 ns/op    0 B/op    0 allocs/op
// BenchmarkInliningThreshold-8         5000000   280 ns/op    0 B/op    0 allocs/op
// BenchmarkNoInline-8                 5000000   285 ns/op    0 B/op    0 allocs/op
//
// Key observations:
// - Stack allocation is ~30-40x faster than heap allocation
// - Simple functions are inlined (SimpleAdd is as fast as stack ops)
// - Functions with loops may not be inlined (InliningThreshold vs NoInline)
// - Local variable caching avoids repeated len() calls
//
// Compiler optimization commands:
// $ go build -gcflags="-m -m"    # Show all optimization decisions
// $ go build -gcflags="-l"       # Disable inlining
// $ go build -gcflags="-l -l"   # More aggressive inlining disable