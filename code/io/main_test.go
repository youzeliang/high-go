package ioopt

import (
	"bufio"
	"bytes"
	"io"
	"testing"
)

// =============================================================================
// bufio benchmarks
// =============================================================================

// BenchmarkBufioReaderRead benchmarks bufio.Reader Read
func BenchmarkBufioReaderRead(b *testing.B) {
	data := bytes.Repeat([]byte("hello world "), 1000)
	reader := bufio.NewReaderSize(bytes.NewReader(data), 4096)
	buf := make([]byte, 4096)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader.Reset(bytes.NewReader(data))
		for {
			_, err := reader.Read(buf)
			if err != nil {
				break
			}
		}
	}
}

// BenchmarkNativeRead benchmarks native io.Reader Read
func BenchmarkNativeRead(b *testing.B) {
	data := bytes.Repeat([]byte("hello world "), 1000)
	buf := make([]byte, 4096)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(data)
		for {
			_, err := reader.Read(buf)
			if err != nil {
				break
			}
		}
	}
}

// BenchmarkBufioWriterWrite benchmarks bufio.Writer Write
func BenchmarkBufioWriterWrite(b *testing.B) {
	data := bytes.Repeat([]byte("hello world "), 1000)
	var buf bytes.Buffer
	writer := bufio.NewWriterSize(&buf, 4096)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		writer.Reset(&buf)
		writer.Write(data)
		writer.Flush()
	}
}

// BenchmarkNativeWriterWrite benchmarks native Write
func BenchmarkNativeWriterWrite(b *testing.B) {
	data := bytes.Repeat([]byte("hello world "), 1000)
	var buf bytes.Buffer
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		buf.Write(data)
	}
}

// =============================================================================
// bytes.Buffer benchmarks
// =============================================================================

// BenchmarkBufferWrite benchmarks bytes.Buffer WriteString
func BenchmarkBufferWrite(b *testing.B) {
	var buf bytes.Buffer
	data := "hello world "
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		for j := 0; j < 1000; j++ {
			buf.WriteString(data)
		}
	}
}

// BenchmarkPreallocBufferWrite benchmarks pre-allocated bytes.Buffer
func BenchmarkPreallocBufferWrite(b *testing.B) {
	data := "hello world "
	size := 1000 * len(data)
	buf := bytes.NewBuffer(make([]byte, 0, size))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		for j := 0; j < 1000; j++ {
			buf.WriteString(data)
		}
	}
}

// BenchmarkBufferWriteInt benchmarks writing int to buffer
func BenchmarkBufferWriteInt(b *testing.B) {
	var buf bytes.Buffer
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		for j := 0; j < 1000; j++ {
			buf.WriteString(string(rune(j)))
		}
	}
}

// =============================================================================
// ReadAll alternatives benchmarks
// =============================================================================

// BenchmarkReadAllNaive benchmarks naive ReadAll
func BenchmarkReadAllNaive(b *testing.B) {
	data := bytes.Repeat([]byte("hello world "), 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ReadAllNaive(bytes.NewReader(data))
	}
}

// BenchmarkReadAllWithBufio benchmarks ReadAll with bufio
func BenchmarkReadAllWithBufio(b *testing.B) {
	data := bytes.Repeat([]byte("hello world "), 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ReadAllWithBufio(bytes.NewReader(data))
	}
}

// BenchmarkReadAllWithBuffer benchmarks ReadAll with pre-allocated buffer
func BenchmarkReadAllWithBuffer(b *testing.B) {
	data := bytes.Repeat([]byte("hello world "), 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ReadAllWithBuffer(bytes.NewReader(data), 4096)
	}
}

// BenchmarkIoCopy benchmarks io.Copy
func BenchmarkIoCopy(b *testing.B) {
	data := bytes.Repeat([]byte("hello world "), 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, bytes.NewReader(data))
	}
}

// =============================================================================
// String building benchmarks
// =============================================================================

// BenchmarkBuildStringNaive benchmarks naive string building
func BenchmarkBuildStringNaive(b *testing.B) {
	parts := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		parts[i] = "hello world "
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BuildStringNaive(parts)
	}
}

// BenchmarkBuildStringWithBuilder benchmarks string building with strings.Builder
func BenchmarkBuildStringWithBuilder(b *testing.B) {
	parts := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		parts[i] = "hello world "
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BuildStringWithBuilder(parts)
	}
}

// BenchmarkBuildStringWithBuffer benchmarks string building with bytes.Buffer
func BenchmarkBuildStringWithBuffer(b *testing.B) {
	parts := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		parts[i] = "hello world "
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BuildStringWithBuffer(parts)
	}
}

// BenchmarkBuildStringWithJoin benchmarks string building with strings.Join
func BenchmarkBuildStringWithJoin(b *testing.B) {
	parts := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		parts[i] = "hello world "
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BuildStringWithJoin(parts)
	}
}

// Example output:
// $ go test -bench="Bufio|Buffer|ReadAll|StringBuild" -benchmem ./code/io/
// BenchmarkBufioReaderRead-8              100000    15000 ns/op   8192 B/op    1 allocs/op
// BenchmarkNativeRead-8                    50000    32000 ns/op  16384 B/op    2 allocs/op
// BenchmarkBufferWrite-8                   500000     3200 ns/op   8000 B/op    1 allocs/op
// BenchmarkPreallocBufferWrite-8          1000000     1200 ns/op   4096 B/op    1 allocs/op
// BenchmarkReadAllNaive-8                  50000    25000 ns/op  12000 B/op    5 allocs/op
// BenchmarkReadAllWithBufio-8            100000    12000 ns/op   4096 B/op    1 allocs/op
// BenchmarkBuildStringNaive-8              10000   150000 ns/op 500000 B/op  1000 allocs/op
// BenchmarkBuildStringWithBuilder-8       500000     3000 ns/op   8192 B/op    1 allocs/op
// BenchmarkBuildStringWithJoin-8         2000000      800 ns/op   4096 B/op    1 allocs/op
