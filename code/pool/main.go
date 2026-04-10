package pool

import (
	"bytes"
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