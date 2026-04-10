package pool

import (
	"bytes"
	"encoding/json"
	"sync"
)

// Pool represents a sync.Pool example for reusing objects
type Pool struct {
	slice sync.Pool
}

// NewPool creates a new Pool
func NewPool() *Pool {
	return &Pool{
		slice: sync.Pool{
			New: func() interface{} {
				return make([]byte, 1024)
			},
		},
	}
}

// Get retrieves a byte slice from the pool
func (p *Pool) Get() []byte {
	return p.slice.Get().([]byte)
}

// Put returns a byte slice to the pool
func (p *Pool) Put(b []byte) {
	p.slice.Put(b)
}

// BufferPool provides pooled bytes.Buffer instances
type BufferPool struct {
	pool sync.Pool
}

// NewBufferPool creates a new BufferPool
func NewBufferPool() *BufferPool {
	return &BufferPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &bytes.Buffer{}
			},
		},
	}
}

// Get retrieves a buffer from the pool
func (bp *BufferPool) Get() *bytes.Buffer {
	buf := bp.pool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

// Put returns a buffer to the pool
func (bp *BufferPool) Put(buf *bytes.Buffer) {
	bp.pool.Put(buf)
}

// ObjectPool is a generic pool for any type
type ObjectPool struct {
	pool sync.Pool
}

// NewObjectPool creates a new ObjectPool with the given factory
func NewObjectPool(factory func() interface{}) *ObjectPool {
	return &ObjectPool{
		pool: sync.Pool{
			New: factory,
		},
	}
}

// Get retrieves an object from the pool
func (op *ObjectPool) Get() interface{} {
	return op.pool.Get()
}

// Put returns an object to the pool
func (op *ObjectPool) Put(obj interface{}) {
	op.pool.Put(obj)
}

// =============================================================================
// Advanced: GC Pressure Reduction Techniques
// =============================================================================

// JSONBufferPool provides pooled bytes.Buffer for JSON encoding
// This reduces GC pressure by reusing buffers instead of allocating new ones
type JSONBufferPool struct {
	pool sync.Pool
}

// NewJSONBufferPool creates a new JSONBufferPool
func NewJSONBufferPool() *JSONBufferPool {
	return &JSONBufferPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &bytes.Buffer{}
			},
		},
	}
}

// Get retrieves a buffer from the pool
func (p *JSONBufferPool) Get() *bytes.Buffer {
	buf := p.pool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

// Put returns a buffer to the pool
func (p *JSONBufferPool) Put(buf *bytes.Buffer) {
	p.pool.Put(buf)
}

// EncodeJSON encodes a value to JSON using pooled buffer
func (p *JSONBufferPool) EncodeJSON(v interface{}) ([]byte, error) {
	buf := p.Get()
	defer p.Put(buf)
	err := json.NewEncoder(buf).Encode(v)
	if err != nil {
		return nil, err
	}
	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())
	return result, nil
}

// DecodeJSON decodes JSON from bytes using pooled buffer
func (p *JSONBufferPool) DecodeJSON(data []byte, v interface{}) error {
	buf := p.Get()
	defer p.Put(buf)
	buf.Write(data)
	return json.NewDecoder(buf).Decode(v)
}

// =============================================================================
// MultiSizeObjectPool - Pool for objects of different sizes
// Reduces GC pressure by maintaining size-specific pools
// =============================================================================

// MultiSizeObjectPool maintains separate pools for different size classes
type MultiSizeObjectPool struct {
	pools []sync.Pool
	sizes []int
}

// NewMultiSizeObjectPool creates a pool for common buffer sizes
func NewMultiSizeObjectPool() *MultiSizeObjectPool {
	sizes := []int{64, 128, 256, 512, 1024, 2048, 4096, 8192}
	pools := make([]sync.Pool, len(sizes))
	for i := range pools {
		size := sizes[i]
		pools[i].New = func() interface{} {
			return make([]byte, size)
		}
	}
	return &MultiSizeObjectPool{
		pools: pools,
		sizes: sizes,
	}
}

// Get retrieves a buffer of at least the requested size
func (m *MultiSizeObjectPool) Get(minSize int) []byte {
	for i, size := range m.sizes {
		if size >= minSize {
			return m.pools[i].Get().([]byte)
		}
	}
	// Fallback: allocate exact size (beyond max size)
	return make([]byte, minSize)
}

// Put returns a buffer to the appropriate pool
func (m *MultiSizeObjectPool) Put(b []byte) {
	for i, size := range m.sizes {
		if size == cap(b) {
			m.pools[i].Put(b)
			return
		}
	}
	// Don't pool non-standard sizes
}

// =============================================================================
// RowBufferPool - Real-world example for database row buffering
// Demonstrates significant GC reduction in data processing
// =============================================================================

// RowBuffer is used to buffer database rows
type RowBuffer struct {
	Columns []string
	Values  []interface{}
}

// RowBufferPool provides pooled RowBuffer instances
type RowBufferPool struct {
	pool sync.Pool
}

// NewRowBufferPool creates a new RowBufferPool
func NewRowBufferPool() *RowBufferPool {
	return &RowBufferPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &RowBuffer{
					Columns: make([]string, 0, 16),
					Values:  make([]interface{}, 0, 16),
				}
			},
		},
	}
}

// Get retrieves a RowBuffer from the pool
func (rbp *RowBufferPool) Get() *RowBuffer {
	rb := rbp.pool.Get().(*RowBuffer)
	rb.Columns = rb.Columns[:0]
	rb.Values = rb.Values[:0]
	return rb
}

// Put returns a RowBuffer to the pool
func (rbp *RowBufferPool) Put(rb *RowBuffer) {
	rbp.pool.Put(rb)
}