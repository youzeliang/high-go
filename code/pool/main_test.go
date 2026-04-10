package pool

import (
	"bytes"
	"testing"
)

// BenchmarkNoPool allocates new buffers each time without pooling
func BenchmarkNoPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := make([]byte, 1024)
		for j := 0; j < 1024; j++ {
			buf[j] = byte(j % 256)
		}
		_ = buf
	}
}

// BenchmarkWithPool uses sync.Pool to reuse byte slices
func BenchmarkWithPool(b *testing.B) {
	pool := NewPool()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := pool.Get()
		for j := 0; j < 1024; j++ {
			buf[j] = byte(j % 256)
		}
		pool.Put(buf)
	}
}

// BenchmarkNoBufferPool allocates new bytes.Buffer each time
func BenchmarkNoBufferPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := new(bytes.Buffer)
		for j := 0; j < 1024; j++ {
			buf.WriteByte(byte(j % 256))
		}
		_ = buf.String()
	}
}

// BenchmarkWithBufferPool uses BufferPool to reuse buffers
func BenchmarkWithBufferPool(b *testing.B) {
	pool := NewBufferPool()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := pool.Get()
		for j := 0; j < 1024; j++ {
			buf.WriteByte(byte(j % 256))
		}
		_ = buf.String()
		pool.Put(buf)
	}
}

// BenchmarkObjectPool demonstrates generic object pooling
func BenchmarkObjectPool(b *testing.B) {
	pool := NewObjectPool(func() interface{} {
		return make([]int, 100)
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		arr := pool.Get().([]int)
		for j := 0; j < 100; j++ {
			arr[j] = j
		}
		pool.Put(arr)
	}
}

// BenchmarkNoObjectPool allocates new slices each time
func BenchmarkNoObjectPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		arr := make([]int, 100)
		for j := 0; j < 100; j++ {
			arr[j] = j
		}
		_ = arr
	}
}

// Example output:
// $ go test -bench="Pool$" -benchmem ./code/pool/
// BenchmarkNoPool-8             500000              3284 ns/op            4096 B/op        1 allocs/op
// BenchmarkWithPool-8           5000000               286 ns/op              0 B/op        0 allocs/op
// BenchmarkNoBufferPool-8        200000              6782 ns/op            4096 B/op        1 allocs/op
// BenchmarkWithBufferPool-8    2000000               812 ns/op              0 B/op        0 allocs/op
// BenchmarkObjectPool-8         3000000               458 ns/op              0 B/op        0 allocs/op
// BenchmarkNoObjectPool-8       1000000              1456 ns/op            800 B/op        1 allocs/op