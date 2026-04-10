package channel

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// =============================================================================
// Unbuffered Channel - Blocks until both sender and receiver are ready
// =============================================================================

// UnbufferedChannel demonstrates a channel with no buffer
// - Send blocks until receiver is ready
// - Receive blocks until sender is ready
// - Ensures synchronous handoff between goroutines
type UnbufferedChannel struct{}

// SendUnbuffered sends a value on an unbuffered channel
// This will block until a receiver is ready
func (u *UnbufferedChannel) Send(ch chan int, value int) {
	ch <- value // Blocks until receiver reads
}

// ReceiveUnbuffered receives a value from an unbuffered channel
// This will block until a sender is ready
func (u *UnbufferedChannel) Receive(ch chan int) int {
	return <-ch // Blocks until sender sends
}

// NewUnbuffered creates an unbuffered channel
func NewUnbuffered() chan int {
	return make(chan int)
}

// =============================================================================
// Buffered Channel - Non-blocking send/receive until buffer is full/empty
// =============================================================================

// BufferedChannel demonstrates a channel with fixed buffer capacity
// - Send blocks only when buffer is full
// - Receive blocks only when buffer is empty
// - Allows asynchronous communication up to buffer capacity
type BufferedChannel struct {
	bufferSize int
}

// NewBuffered creates a buffered channel with specified size
func NewBuffered(size int) chan int {
	return make(chan int, size)
}

// SendBuffered sends a value on a buffered channel
// Non-blocking if buffer has space
func (b *BufferedChannel) Send(ch chan int, value int) {
	ch <- value // Non-blocking if buffer has space
}

// ReceiveBuffered receives a value from a buffered channel
// Non-blocking if buffer has data
func (b *BufferedChannel) Receive(ch chan int) int {
	return <-ch // Non-blocking if buffer has data
}

// =============================================================================
// Channel Select - Non-blocking multi-channel operations
// =============================================================================

// SelectDemo demonstrates using select for non-blocking channel operations
type SelectDemo struct{}

// SelectReceive uses select to receive from multiple channels
// Returns immediately if no channel has data
func (s *SelectDemo) SelectReceive(ch1, ch2 <-chan int) (int, bool) {
	select {
	case v := <-ch1:
		return v, true
	case v := <-ch2:
		return v, true
	default:
		return 0, false // No data available, non-blocking
	}
}

// SelectSend uses select to send to a channel without blocking
func (s *SelectDemo) SelectSend(ch chan int, value int) bool {
	select {
	case ch <- value:
		return true
	default:
		return false // Buffer full, non-blocking
	}
}

// SelectTimeout demonstrates select with timeout
func (s *SelectDemo) SelectTimeout(ch <-chan int, timeout time.Duration) (int, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	select {
	case v := <-ch:
		return v, true
	case <-ctx.Done():
		return 0, false // Timeout
	}
}

// =============================================================================
// Channel Patterns
// =============================================================================

// FanOut demonstrates distributing work across multiple workers
func FanOut(workers int, jobs <-chan int, results chan<- int) {
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for job := range jobs {
				// Simulate work
				time.Sleep(time.Microsecond * 10)
				results <- job * 2
			}
		}(i)
	}
	wg.Wait()
	close(results)
}

// Pipeline demonstrates chaining channel operations
func Pipeline(input <-chan int) <-chan int {
	// Stage 1: Double
	out1 := make(chan int)
	go func() {
		for v := range input {
			out1 <- v * 2
		}
		close(out1)
	}()

	// Stage 2: Add 10
	out2 := make(chan int)
	go func() {
		for v := range out1 {
			out2 <- v + 10
		}
		close(out2)
	}()

	return out2
}

// =============================================================================
// Channel Metrics
// =============================================================================

// ChannelStats holds performance metrics for channel operations
type ChannelStats struct {
	SendCount    int
	ReceiveCount int
	SendDuration time.Duration
}

// MeasureUnbuffered measures unbuffered channel performance
func MeasureUnbuffered(iterations int) ChannelStats {
	ch := make(chan int)
	var stats ChannelStats

	start := time.Now()
	for i := 0; i < iterations; i++ {
		go func(v int) {
			ch <- v
		}(i)
		<-ch
	}
	stats.SendDuration = time.Since(start)
	stats.SendCount = iterations
	stats.ReceiveCount = iterations

	return stats
}

// MeasureBuffered measures buffered channel performance
func MeasureBuffered(iterations int, bufferSize int) ChannelStats {
	ch := make(chan int, bufferSize)
	var stats ChannelStats

	start := time.Now()
	for i := 0; i < iterations; i++ {
		ch <- i
	}
	for i := 0; i < iterations; i++ {
		<-ch
	}
	stats.SendDuration = time.Since(start)
	stats.SendCount = iterations
	stats.ReceiveCount = iterations

	return stats
}

// =============================================================================
// Context with Channels
// =============================================================================

// ContextWithChannel demonstrates channel with context cancellation
type ContextWithChannel struct{}

// ProducerWithCancel produces values until context is cancelled
func (c *ContextWithChannel) ProducerWithCancel(ctx context.Context, ch chan<- int) {
	for i := 0; ; i++ {
		select {
		case ch <- i:
			// Value sent successfully
		case <-ctx.Done():
			close(ch)
			return
		}
	}
}

// ConsumerWithTimeout consumes values with timeout
func (c *ContextWithChannel) ConsumerWithTimeout(ctx context.Context, ch <-chan int, timeout time.Duration) ([]int, error) {
	result := make([]int, 0)
	ticker := time.NewTicker(timeout)
	defer ticker.Stop()

	for {
		select {
		case v, ok := <-ch:
			if !ok {
				return result, nil
			}
			result = append(result, v)
		case <-ticker.C:
			return result, fmt.Errorf("timeout after %v", timeout)
		case <-ctx.Done():
			return result, ctx.Err()
		}
	}
}
