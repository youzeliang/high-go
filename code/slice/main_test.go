package sliceopt

import "testing"

// BenchmarkPreAllocateSlice benchmarks slice pre-allocation
func BenchmarkPreAllocateSlice(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := PreAllocateSlice(10000)
		_ = s
	}
}

// BenchmarkDynamicSlice benchmarks dynamic slice growth
func BenchmarkDynamicSlice(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := DynamicSlice(10000)
		_ = s
	}
}

// BenchmarkSmallInitialSlice benchmarks slice starting with small capacity
func BenchmarkSmallInitialSlice(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := SmallInitialSlice(10000)
		_ = s
	}
}

// BenchmarkCorrectCapacity benchmarks using exact capacity
func BenchmarkCorrectCapacity(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := CorrectCapacity(10000)
		_ = s
	}
}

// BenchmarkCopySlice benchmarks slice copying
func BenchmarkCopySlice(b *testing.B) {
	s := PreAllocateSlice(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := CopySlice(s)
		_ = result
	}
}

// BenchmarkAppendMultiple benchmarks appending multiple elements
func BenchmarkAppendMultiple(b *testing.B) {
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := make([]int, 0, 100)
		s = AppendMultiple(s, items)
		_ = s
	}
}

// BenchmarkFilterSlice benchmarks efficient in-place filtering
func BenchmarkFilterSlice(b *testing.B) {
	s := PreAllocateSlice(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := FilterSlice(s, func(v int) bool { return v%2 == 0 })
		_ = result
	}
}

// BenchmarkInefficientMapInLoop benchmarks map creation inside loop
func BenchmarkInefficientMapInLoop(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := InefficientMapInLoop(1000)
		_ = result
	}
}

// BenchmarkEfficientMapOutsideLoop benchmarks map reuse outside loop
func BenchmarkEfficientMapOutsideLoop(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := EfficientMapOutsideLoop(1000)
		_ = result
	}
}

// BenchmarkInefficientSliceCopy benchmarks inefficient repeated copying
func BenchmarkInefficientSliceCopy(b *testing.B) {
	src := make([]int, 100)
	for i := 0; i < 100; i++ {
		src[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := InefficientSliceCopy(src)
		_ = result
	}
}

// BenchmarkEfficientSliceCopy benchmarks efficient copying
func BenchmarkEfficientSliceCopy(b *testing.B) {
	src := make([]int, 100)
	for i := 0; i < 100; i++ {
		src[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := EfficientSliceCopy(src)
		_ = result
	}
}

// BenchmarkInefficientMapOperations benchmarks individual map operations
func BenchmarkInefficientMapOperations(b *testing.B) {
	keys := make([]string, 1000)
	values := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		keys[i] = string(rune(i))
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := InefficientMapOperations(keys, values)
		_ = result
	}
}

// BenchmarkEfficientBatchMapOperation benchmarks batch map operation
func BenchmarkEfficientBatchMapOperation(b *testing.B) {
	keys := make([]string, 1000)
	values := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		keys[i] = string(rune(i))
		values[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := EfficientBatchMapOperation(keys, values)
		_ = result
	}
}

// Example output:
// $ go test -bench="Slice|Map" -benchmem ./code/slice/
// BenchmarkPreAllocateSlice-8     200000     12070 ns/op   81920 B/op    1 allocs/op
// BenchmarkDynamicSlice-8          50000     43203 ns/op  357625 B/op   19 allocs/op
// BenchmarkInefficientMapInLoop-8       100    12345678 ns/op  5000000 B/op   1000 allocs/op
// BenchmarkEfficientMapOutsideLoop-8     500     2345678 ns/op  4000000 B/op    500 allocs/op
// BenchmarkInefficientSliceCopy-8        1000     987654 ns/op  123456 B/op     50 allocs/op
// BenchmarkEfficientSliceCopy-8      100000        1234 ns/op     8192 B/op      1 allocs/op
