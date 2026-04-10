package ioopt

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
)

// =============================================================================
// bufio buffered I/O
// =============================================================================

// ReaderWithBufio demonstrates using bufio.Reader for efficient reading
type ReaderWithBufio struct {
	reader *bufio.Reader
}

// NewReaderWithBufio creates a bufio.Reader
func NewReaderWithBufio(r io.Reader) *ReaderWithBufio {
	return &ReaderWithBufio{
		reader: bufio.NewReaderSize(r, 64*1024), // 64KB buffer
	}
}

// ReadLine reads a line using bufio.Reader
func (r *ReaderWithBufio) ReadLine() ([]byte, error) {
	return r.reader.ReadBytes('\n')
}

// ReadWord reads a word using bufio.Reader
func (r *ReaderWithBufio) ReadWord() (string, error) {
	return r.reader.ReadString(' ')
}

// WriterWithBufio demonstrates using bufio.Writer for efficient writing
type WriterWithBufio struct {
	writer *bufio.Writer
}

// NewWriterWithBufio creates a bufio.Writer
func NewWriterWithBufio(w io.Writer) *WriterWithBufio {
	return &WriterWithBufio{
		writer: bufio.NewWriterSize(w, 64*1024), // 64KB buffer
	}
}

// WriteString writes a string using bufio.Writer
func (w *WriterWithBufio) WriteString(s string) (int, error) {
	return w.writer.WriteString(s)
}

// Flush flushes the writer
func (w *WriterWithBufio) Flush() error {
	return w.writer.Flush()
}

// =============================================================================
// bytes.Buffer operations
// =============================================================================

// BufferConcat uses bytes.Buffer for string concatenation
type BufferConcat struct {
	buf *bytes.Buffer
}

// NewBufferConcat creates a new BufferConcat
func NewBufferConcat() *BufferConcat {
	return &BufferConcat{
		buf: new(bytes.Buffer),
	}
}

// Write writes data to buffer
func (b *BufferConcat) Write(s string) {
	b.buf.WriteString(s)
}

// WriteInt writes an int to buffer
func (b *BufferConcat) WriteInt(i int) {
	b.buf.WriteString(string(rune(i)))
}

// String returns the buffer contents as string
func (b *BufferConcat) String() string {
	return b.buf.String()
}

// Reset resets the buffer
func (b *BufferConcat) Reset() {
	b.buf.Reset()
}

// PreallocBufferConcat uses pre-allocated buffer
type PreallocBufferConcat struct {
	buf *bytes.Buffer
}

// NewPreallocBufferConcat creates a pre-allocated buffer
func NewPreallocBufferConcat(size int) *PreallocBufferConcat {
	return &PreallocBufferConcat{
		buf: bytes.NewBuffer(make([]byte, 0, size)),
	}
}

// Write writes data to pre-allocated buffer
func (b *PreallocBufferConcat) Write(s string) {
	b.buf.WriteString(s)
}

// String returns the buffer contents
func (b *PreallocBufferConcat) String() string {
	return b.buf.String()
}

// Reset resets the buffer
func (b *PreallocBufferConcat) Reset() {
	b.buf.Reset()
}

// =============================================================================
// ReadAll alternatives
// =============================================================================

// ReadAllNaive uses naive approach with multiple allocations
func ReadAllNaive(r io.Reader) ([]byte, error) {
	var result []byte
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			result = append(result, buf[:n]...)
		}
		if err != nil {
			if err == io.EOF {
				return result, nil
			}
			return nil, err
		}
	}
}

// ReadAllWithBufio uses bufio.Reader for efficient reading
func ReadAllWithBufio(r io.Reader) ([]byte, error) {
	reader := bufio.NewReader(r)
	return io.ReadAll(reader)
}

// ReadAllWithBuffer uses pre-allocated buffer
func ReadAllWithBuffer(r io.Reader, bufSize int) ([]byte, error) {
	reader := bufio.NewReaderSize(r, bufSize)
	return io.ReadAll(reader)
}

// ReadAllUsingCopy uses io.Copy to read all
func ReadAllUsingCopy(w io.Writer, r io.Reader) (int64, error) {
	return io.Copy(w, r)
}

// =============================================================================
// File I/O with bufio
// =============================================================================

// ReadFileWithBufio reads entire file using bufio
func ReadFileWithBufio(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReaderSize(file, 1024*1024) // 1MB buffer
	return io.ReadAll(reader)
}

// WriteFileWithBufio writes data to file using bufio
func WriteFileWithBufio(filename string, data []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriterSize(file, 1024*1024) // 1MB buffer
	_, err = writer.Write(data)
	if err != nil {
		return err
	}
	return writer.Flush()
}

// ReadLinesWithBufio reads lines efficiently using bufio.Scanner
func ReadLinesWithBufio(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// ReadLargeFileChunked reads large file in chunks
func ReadLargeFileChunked(filename string, chunkSize int) ([][]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var chunks [][]byte
	buf := make([]byte, chunkSize)
	for {
		n, err := file.Read(buf)
		if n > 0 {
			chunk := make([]byte, n)
			copy(chunk, buf[:n])
			chunks = append(chunks, chunk)
		}
		if err != nil {
			if err == io.EOF {
				return chunks, nil
			}
			return nil, err
		}
	}
}

// =============================================================================
// String building alternatives
// =============================================================================

// BuildStringNaive builds string with naive concatenation
func BuildStringNaive(parts []string) string {
	result := ""
	for _, p := range parts {
		result += p
	}
	return result
}

// BuildStringWithBuilder builds string with strings.Builder
func BuildStringWithBuilder(parts []string) string {
	var builder strings.Builder
	for _, p := range parts {
		builder.WriteString(p)
	}
	return builder.String()
}

// BuildStringWithBuffer builds string with bytes.Buffer
func BuildStringWithBuffer(parts []string) string {
	buf := new(bytes.Buffer)
	for _, p := range parts {
		buf.WriteString(p)
	}
	return buf.String()
}

// BuildStringWithJoin builds string with strings.Join (most efficient)
func BuildStringWithJoin(parts []string) string {
	return strings.Join(parts, "")
}
