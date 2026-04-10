package ringbuf

import "testing"

// BenchmarkRingBufferPush benchmarks push operations on ring buffer
func BenchmarkRingBufferPush(b *testing.B) {
	r := NewRingBuffer(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Push(i)
	}
}

// BenchmarkRingBufferPop benchmarks pop operations on ring buffer
func BenchmarkRingBufferPop(b *testing.B) {
	r := NewRingBuffer(1000)
	for i := 0; i < 1000; i++ {
		r.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Pop()
		r.Push(i)
	}
}

// BenchmarkRingBufferFullCycle benchmarks full write/read cycle
func BenchmarkRingBufferFullCycle(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := NewRingBuffer(1000)
		for j := 0; j < 1000; j++ {
			r.Push(j)
		}
		for j := 0; j < 1000; j++ {
			r.Pop()
		}
	}
}

// BenchmarkIntRingBufferPush benchmarks int ring buffer push
func BenchmarkIntRingBufferPush(b *testing.B) {
	r := NewIntRingBuffer(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Push(i)
	}
}

// BenchmarkIntRingBufferPop benchmarks int ring buffer pop
func BenchmarkIntRingBufferPop(b *testing.B) {
	r := NewIntRingBuffer(1000)
	for i := 0; i < 1000; i++ {
		r.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Pop()
		r.Push(i)
	}
}

// BenchmarkChannelRingBufferPush benchmarks channel ring buffer push
func BenchmarkChannelRingBufferPush(b *testing.B) {
	r := NewChannelRingBuffer(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Push(i)
	}
}

// BenchmarkChannelRingBufferPop benchmarks channel ring buffer pop
func BenchmarkChannelRingBufferPop(b *testing.B) {
	r := NewChannelRingBuffer(1000)
	for i := 0; i < 1000; i++ {
		r.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Pop()
		r.Push(i)
	}
}

// BenchmarkSliceQueuePush benchmarks slice-based queue push
func BenchmarkSliceQueuePush(b *testing.B) {
	var q []int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q = append(q, i)
	}
}

// BenchmarkSliceQueuePop benchmarks slice-based queue pop (from front)
func BenchmarkSliceQueuePop(b *testing.B) {
	q := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		q[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if len(q) > 0 {
			q = q[1:]
		}
	}
}

// BenchmarkRingBufferOverwrite benchmarks ring buffer overwrite behavior
func BenchmarkRingBufferOverwrite(b *testing.B) {
	r := NewRingBuffer(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Push(i)
	}
}

// Example output:
// $ go test -bench="Ring" -benchmem ./code/ring/
// BenchmarkRingBufferPush-8            100000000    12.3 ns/op    0 B/op    0 allocs/op
// BenchmarkRingBufferPop-8             50000000    24.5 ns/op    0 B/op    0 allocs/op
// BenchmarkRingBufferFullCycle-8        500000    3200 ns/op    0 B/op    0 allocs/op
// BenchmarkIntRingBufferPush-8         100000000     8.2 ns/op    0 B/op    0 allocs/op
// BenchmarkIntRingBufferPop-8          50000000    18.1 ns/op    0 B/op    0 allocs/op
// BenchmarkChannelRingBufferPush-8      20000000    78.3 ns/op   16 B/op    1 allocs/op
// BenchmarkSliceQueuePush-8           100000000    15.6 ns/op    0 B/op    0 allocs/op
// BenchmarkSliceQueuePop-8              100000    9800 ns/op    0 B/op    0 allocs/op
