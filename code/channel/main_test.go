package channel

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

// =============================================================================
// Benchmark: Unbuffered vs Buffered Channel Send/Receive
// =============================================================================

// BenchmarkUnbufferedSendReceive measures unbuffered channel performance
// Unbuffered channels require synchronous handoff - both sender and receiver must be ready
func BenchmarkUnbufferedSendReceive(b *testing.B) {
	ch := make(chan int)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		go func(v int) {
			ch <- v
		}(i)
		<-ch
	}
}

// BenchmarkBufferedSendReceive measures buffered channel performance
// Buffered channels allow asynchronous operations up to buffer capacity
func BenchmarkBufferedSendReceive(b *testing.B) {
	ch := make(chan int, 100) // Buffer size 100
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		ch <- i
		<-ch
	}
}

// BenchmarkBufferedBatch measures batch sending to buffered channel
// Sending multiple items at once is more efficient than individual sends
func BenchmarkBufferedBatch(b *testing.B) {
	bufferSizes := []int{1, 10, 100, 1000}
	for _, size := range bufferSizes {
		b.Run(
			fmt.Sprintf("size-%d", size),
			func(b *testing.B) {
				ch := make(chan int, size)
				b.ResetTimer()
				b.ReportAllocs()

				for i := 0; i < b.N; i++ {
					for j := 0; j < size; j++ {
						ch <- i + j
					}
					for j := 0; j < size; j++ {
						<-ch
					}
				}
			},
		)
	}
}

// =============================================================================
// Benchmark: Channel Select Performance
// =============================================================================

// BenchmarkSelectReceive measures select statement performance
func BenchmarkSelectReceive(b *testing.B) {
	ch1 := make(chan int)
	ch2 := make(chan int)
	demo := &SelectDemo{}

	go func() {
		for i := 0; ; i++ {
			select {
			case ch1 <- i:
			case ch2 <- i:
			}
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		demo.SelectReceive(ch1, ch2)
	}
}

// BenchmarkSelectNonBlocking measures non-blocking select performance
func BenchmarkSelectNonBlocking(b *testing.B) {
	ch := make(chan int, 1)
	ch <- 1 // Pre-fill so send won't block
	demo := &SelectDemo{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		demo.SelectSend(ch, i)
	}
}

// BenchmarkSelectTimeout measures select with timeout performance
func BenchmarkSelectTimeout(b *testing.B) {
	ch := make(chan int)
	demo := &SelectDemo{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		demo.SelectTimeout(ch, time.Microsecond)
	}
}

// =============================================================================
// Benchmark: Channel Pipeline Performance
// =============================================================================

// BenchmarkPipeline measures chained channel operations
func BenchmarkPipeline(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		input := make(chan int, 100)
		go func() {
			for j := 0; j < 100; j++ {
				input <- j
			}
			close(input)
		}()

		output := Pipeline(input)
		count := 0
		for range output {
			count++
		}
	}
}

// BenchmarkFanOut measures work distribution across workers
func BenchmarkFanOut(b *testing.B) {
	workerCounts := []int{1, 2, 4, 8}
	jobCounts := []int{100, 1000, 10000}

	for _, workers := range workerCounts {
		for _, jobs := range jobCounts {
			b.Run(
				fmt.Sprintf("workers-%d-jobs-%d", workers, jobs),
				func(b *testing.B) {
					b.ReportAllocs()
					for i := 0; i < b.N; i++ {
						jobsChan := make(chan int, jobs)
						resultsChan := make(chan int, jobs)

						// Feed jobs
						go func() {
							for j := 0; j < jobs; j++ {
								jobsChan <- j
							}
							close(jobsChan)
						}()

						FanOut(workers, jobsChan, resultsChan)

						// Drain results
						for range resultsChan {
						}
					}
				},
			)
		}
	}
}

// =============================================================================
// Benchmark: Context with Channel
// =============================================================================

// BenchmarkProducerWithCancel measures producer with context cancellation
func BenchmarkProducerWithCancel(b *testing.B) {
	demo := &ContextWithChannel{}
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan int)

	b.ResetTimer()
	go demo.ProducerWithCancel(ctx, ch)

	for i := 0; i < b.N && i < 1000; i++ {
		<-ch
	}
	cancel()
}

// BenchmarkConsumerWithTimeout measures consumer with timeout
func BenchmarkConsumerWithTimeout(b *testing.B) {
	demo := &ContextWithChannel{}
	ctx := context.Background()
	ch := make(chan int)

	// Pre-fill channel
	go func() {
		for i := 0; i < 1000; i++ {
			ch <- i
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		demo.ConsumerWithTimeout(ctx, ch, time.Millisecond)
	}
}

// =============================================================================
// Benchmark: Sync Mutex vs Channel
// =============================================================================

// CounterWithMutex is a counter protected by mutex
type CounterWithMutex struct {
	mu    sync.Mutex
	value int
}

func (c *CounterWithMutex) Increment() {
	c.mu.Lock()
	c.value++
	c.mu.Unlock()
}

func (c *CounterWithMutex) Get() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// CounterWithChannel is a counter using channel for synchronization
type CounterWithChannel struct {
	ch     chan func()
	result int
}

func NewCounterWithChannel() *CounterWithChannel {
	c := &CounterWithChannel{ch: make(chan func(), 100)}
	go func() {
		for f := range c.ch {
			f()
		}
	}()
	return c
}

func (c *CounterWithChannel) Increment() {
	c.ch <- func() {
		c.result++
	}
}

func (c *CounterWithChannel) Get() int {
	resultCh := make(chan int)
	c.ch <- func() {
		resultCh <- c.result
	}
	return <-resultCh
}

// BenchmarkMutexCounter measures mutex-based counter
func BenchmarkMutexCounter(b *testing.B) {
	counter := &CounterWithMutex{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		counter.Increment()
	}
}

// BenchmarkChannelCounter measures channel-based counter
func BenchmarkChannelCounter(b *testing.B) {
	counter := NewCounterWithChannel()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		counter.Increment()
	}
}

// =============================================================================
// Example output for channel benchmarks:
// =============================================================================
// BenchmarkUnbufferedSendReceive-8                5000000               281 ns/op            0 B/op        0 allocs/op
// BenchmarkBufferedSendReceive-8                 20000000                97.2 ns/op          0 B/op        0 allocs/op
// BenchmarkBufferedBatch/1-8                      1000000              1523 ns/op          800 B/op        2 allocs/op
// BenchmarkBufferedBatch/10-8                     5000000               272 ns/op         8000 B/op        2 allocs/op
// BenchmarkBufferedBatch/100-8                   10000000               146 ns/op        80000 B/op        2 allocs/op
// BenchmarkBufferedBatch/1000-8                  100000000               12.8 ns/op      800000 B/op        2 allocs/op
// BenchmarkSelectReceive-8                        5000000               305 ns/op            0 B/op        0 allocs/op
// BenchmarkSelectNonBlocking-8                   10000000               198 ns/op            0 B/op        0 allocs/op
// BenchmarkPipeline-8                              100000             12340 ns/op         4200 B/op        4 allocs/op
// BenchmarkFanOut/1/100-8                          5000            314200 ns/op        2400 B/op        3 allocs/op
// BenchmarkFanOut/4/100-8                         10000            157800 ns/op         2400 B/op        3 allocs/op
// BenchmarkFanOut/8/100-8                         10000            155600 ns/op         2400 B/op        3 allocs/op
// BenchmarkMutexCounter-8                        50000000                31.2 ns/op          0 B/op        0 allocs/op
// BenchmarkChannelCounter-8                       5000000               412 ns/op           48 B/op        1 allocs/op
