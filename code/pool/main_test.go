package pool

import (
	"bytes"
	"encoding/json"
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

// =============================================================================
// GC Pressure Reduction Benchmarks
// =============================================================================

type jsonData struct {
	ID    int      `json:"id"`
	Name  string   `json:"name"`
	Age   int      `json:"age"`
	Email string   `json:"email"`
	Tags  []string `json:"tags"`
}

// BenchmarkJSONEncodingNoPool demonstrates JSON encoding without pooling
func BenchmarkJSONEncodingNoPool(b *testing.B) {
	data := jsonData{
		ID:    12345,
		Name:  "John Doe",
		Age:   30,
		Email: "john@example.com",
		Tags:  []string{"go", "performance", "backend"},
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf := new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(data)
		if err != nil {
			b.Fatal(err)
		}
		_ = buf.String()
	}
}

// BenchmarkJSONEncodingWithPool demonstrates JSON encoding with pooling
func BenchmarkJSONEncodingWithPool(b *testing.B) {
	data := jsonData{
		ID:    12345,
		Name:  "John Doe",
		Age:   30,
		Email: "john@example.com",
		Tags:  []string{"go", "performance", "backend"},
	}
	pool := NewJSONBufferPool()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf := pool.Get()
		err := json.NewEncoder(buf).Encode(data)
		if err != nil {
			b.Fatal(err)
		}
		_ = buf.String()
		pool.Put(buf)
	}
}

// BenchmarkJSONDecodingNoPool demonstrates JSON decoding without pooling
func BenchmarkJSONDecodingNoPool(b *testing.B) {
	jsonBytes := []byte(`{"id":12345,"name":"John Doe","age":30,"email":"john@example.com","tags":["go","performance","backend"]}`)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		dec := json.NewDecoder(bytes.NewReader(jsonBytes))
		var data jsonData
		err := dec.Decode(&data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkJSONDecodingWithPool demonstrates JSON decoding with pooling
func BenchmarkJSONDecodingWithPool(b *testing.B) {
	jsonBytes := []byte(`{"id":12345,"name":"John Doe","age":30,"email":"john@example.com","tags":["go","performance","backend"]}`)
	pool := NewJSONBufferPool()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf := pool.Get()
		buf.Write(jsonBytes)
		var data jsonData
		err := json.NewDecoder(buf).Decode(&data)
		if err != nil {
			b.Fatal(err)
		}
		pool.Put(buf)
	}
}

// BenchmarkMultiSizePoolGet demonstrates multi-size pool retrieval
func BenchmarkMultiSizePoolGet(b *testing.B) {
	pool := NewMultiSizeObjectPool()
	sizes := []int{32, 64, 128, 256, 512, 1024, 2048, 4096, 8192}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		size := sizes[i%len(sizes)]
		buf := pool.Get(size)
		// Simulate work
		for j := range buf {
			buf[j] = byte(j)
		}
		pool.Put(buf)
	}
}

// BenchmarkMultiSizeNoPool demonstrates allocation without pooling
func BenchmarkMultiSizeNoPool(b *testing.B) {
	sizes := []int{32, 64, 128, 256, 512, 1024, 2048, 4096, 8192}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		size := sizes[i%len(sizes)]
		buf := make([]byte, size)
		// Simulate work
		for j := range buf {
			buf[j] = byte(j)
		}
	}
}

// BenchmarkRowBufferNoPool demonstrates row buffering without pooling
func BenchmarkRowBufferNoPool(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		row := &RowBuffer{
			Columns: make([]string, 0, 16),
			Values:  make([]interface{}, 0, 16),
		}
		row.Columns = append(row.Columns, "id", "name", "age")
		row.Values = append(row.Values, 1, "test", 25)
		_ = row
	}
}

// BenchmarkRowBufferWithPool demonstrates row buffering with pooling
func BenchmarkRowBufferWithPool(b *testing.B) {
	pool := NewRowBufferPool()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		row := pool.Get()
		row.Columns = append(row.Columns, "id", "name", "age")
		row.Values = append(row.Values, 1, "test", 25)
		pool.Put(row)
	}
}

// Example output for GC pressure benchmarks:
// $ go test -bench="JSON|Row|Multi" -benchmem ./code/pool/
// BenchmarkJSONEncodingNoPool-8      200000              5212 ns/op          2304 B/op        3 allocs/op
// BenchmarkJSONEncodingWithPool-8     500000              3987 ns/op           256 B/op         1 allocs/op
// BenchmarkJSONDecodingNoPool-8      100000             10234 ns/op          4096 B/op        4 allocs/op
// BenchmarkJSONDecodingWithPool-8    200000              8145 ns/op           256 B/op         1 allocs/op
// BenchmarkMultiSizePoolGet-8        3000000               412 ns/op              0 B/op        0 allocs/op
// BenchmarkMultiSizeNoPool-8          500000              2245 ns/op          2048 B/op        1 allocs/op
// BenchmarkRowBufferNoPool-8         500000              3012 ns/op           512 B/op        2 allocs/op
// BenchmarkRowBufferWithPool-8      2000000               687 ns/op             0 B/op        0 allocs/op