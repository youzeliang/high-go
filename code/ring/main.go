package ringbuf

import "sync/atomic"

// RingBuffer is a fixed-size circular buffer
type RingBuffer struct {
	buffer    []interface{}
	size      int
	head      int64 // read position
	tail      int64 // write position
	readCnt   int64
	writeCnt  int64
}

// NewRingBuffer creates a new ring buffer with the given size
func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		buffer: make([]interface{}, size),
		size:   size,
	}
}

// Push adds an item to the buffer (overwrites if full)
func (r *RingBuffer) Push(item interface{}) {
	pos := atomic.AddInt64(&r.tail, 1) - 1
	r.buffer[pos%int64(r.size)] = item
	atomic.AddInt64(&r.writeCnt, 1)

	// Handle wrap-around - move head forward if we're overwriting
	if r.writeCnt > int64(r.size) {
		atomic.StoreInt64(&r.head, r.tail-int64(r.size))
	}
}

// Pop removes and returns an item from the buffer
func (r *RingBuffer) Pop() (interface{}, bool) {
	if r.writeCnt <= int64(r.size) {
		// Buffer not yet full
		if r.readCnt >= r.writeCnt {
			return nil, false
		}
	} else {
		// Buffer is full, check if empty due to overread
		if r.tail-r.head == 0 {
			return nil, false
		}
	}

	pos := atomic.AddInt64(&r.head, 1) - 1
	atomic.AddInt64(&r.readCnt, 1)
	return r.buffer[pos%int64(r.size)], true
}

// Len returns the number of items in the buffer
func (r *RingBuffer) Len() int {
	if r.writeCnt <= int64(r.size) {
		return int(r.writeCnt - r.readCnt)
	}
	return int(int64(r.size) - (r.tail - r.head) + int64(r.size))
}

// IsEmpty returns true if the buffer is empty
func (r *RingBuffer) IsEmpty() bool {
	if r.writeCnt <= int64(r.size) {
		return r.readCnt >= r.writeCnt
	}
	return r.tail-r.head == 0
}

// IsFull returns true if the buffer is full
func (r *RingBuffer) IsFull() bool {
	return r.tail-r.head >= int64(r.size)
}

// Capacity returns the buffer capacity
func (r *RingBuffer) Capacity() int {
	return r.size
}

// IntRingBuffer is a type-safe ring buffer for integers
type IntRingBuffer struct {
	buffer []int
	size   int
	head   int64
	tail   int64
}

// NewIntRingBuffer creates a new int ring buffer
func NewIntRingBuffer(size int) *IntRingBuffer {
	return &IntRingBuffer{
		buffer: make([]int, size),
		size:   size,
	}
}

// Push adds an int to the buffer
func (r *IntRingBuffer) Push(item int) {
	pos := atomic.AddInt64(&r.tail, 1) - 1
	r.buffer[pos%int64(r.size)] = item
}

// Pop removes and returns an int from the buffer
func (r *IntRingBuffer) Pop() (int, bool) {
	if r.tail-r.head == 0 {
		return 0, false
	}
	pos := atomic.AddInt64(&r.head, 1) - 1
	return r.buffer[pos%int64(r.size)], true
}

// Len returns the number of items
func (r *IntRingBuffer) Len() int {
	return int(r.tail - r.head)
}

// IsEmpty returns true if empty
func (r *IntRingBuffer) IsEmpty() bool {
	return r.tail == r.head
}

// IsFull returns true if full
func (r *IntRingBuffer) IsFull() bool {
	return r.tail-r.head >= int64(r.size)
}

// Capacity returns the buffer capacity
func (r *IntRingBuffer) Capacity() int {
	return r.size
}

// ChannelRingBuffer uses a channel to implement a ring buffer pattern
type ChannelRingBuffer struct {
	ch   chan interface{}
	size int
}

// NewChannelRingBuffer creates a channel-based ring buffer
func NewChannelRingBuffer(size int) *ChannelRingBuffer {
	return &ChannelRingBuffer{
		ch:   make(chan interface{}, size),
		size: size,
	}
}

// Push adds an item to the buffer (non-blocking)
func (r *ChannelRingBuffer) Push(item interface{}) bool {
	select {
	case r.ch <- item:
		return true
	default:
		return false
	}
}

// Pop removes and returns an item (non-blocking)
func (r *ChannelRingBuffer) Pop() (interface{}, bool) {
	select {
	case item := <-r.ch:
		return item, true
	default:
		return nil, false
	}
}

// Len returns the number of items in the buffer
func (r *ChannelRingBuffer) Len() int {
	return len(r.ch)
}

// IsEmpty returns true if empty
func (r *ChannelRingBuffer) IsEmpty() bool {
	return len(r.ch) == 0
}

// IsFull returns true if full
func (r *ChannelRingBuffer) IsFull() bool {
	return len(r.ch) == r.size
}

// Capacity returns the buffer capacity
func (r *ChannelRingBuffer) Capacity() int {
	return r.size
}
