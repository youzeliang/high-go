# Synchronization Primitives

## 概述

Go 提供了多种同步原语，选择合适的原语对性能至关重要。本章节对比 Mutex、RWMutex、Atomic 和 sync.Once 的性能特点。

## 核心模式

### 1. Mutex vs RWMutex

**Mutex** - 互斥锁，任何时候只允许一个 goroutine 持有锁。

```go
var mu sync.Mutex
mu.Lock()
counter++
mu.Unlock()
```

**RWMutex** - 读写互斥锁，允许多个读锁并行，但写锁独占。

```go
var mu sync.RWMutex
// 读操作
mu.RLock()
value := data
mu.RUnlock()
// 写操作
mu.Lock()
data = newValue
mu.Unlock()
```

### 2. 何时使用 RWMutex

- **读多写少** - 多个 goroutine 同时读取，写入较少
- **读操作耗时** - 读操作本身需要较长时间
- **数据结构较大** - 复制成本高

### 3. Atomic Operations

原子操作是无锁编程的基础，比锁更快。

```go
var counter int64

// 原子增加
atomic.AddInt64(&counter, 1)

// 原子读取
value := atomic.LoadInt64(&counter)

// 原子写入
atomic.StoreInt64(&counter, 42)

// 原子交换
old := atomic.SwapInt64(&counter, 42)

// CAS 操作
atomic.CompareAndSwapInt64(&counter, old, new)
```

### 4. sync.Once

保证expensive initialization 只执行一次。

```go
var once sync.Once
once.Do(func() {
    // 只执行一次
    initExpensiveResource()
})
```

### 5. sync.Cond - 条件变量

`sync.Cond` 用于goroutine 之间的信号通知。

```go
cond := sync.NewCond(&sync.Mutex{})

// 等待条件
cond.L.Lock()
for !condition {
    cond.Wait()
}
cond.L.Unlock()

// 通知一个等待者
cond.Signal()

// 通知所有等待者
cond.Broadcast()
```

**Signal vs Broadcast：**
- `Signal()` - 唤醒一个等待的 goroutine
- `Broadcast()` - 唤醒所有等待的 goroutine

**使用场景：**
- 消费者-生产者模式
- 等待特定条件完成
- 事件通知

## 性能对比

```
BenchmarkMutexWriteRead-8            10047669    118.2 ns/op    0 B/op    0 allocs/op
BenchmarkRWMutexWriteRead-8          19839423     65.69 ns/op    0 B/op    0 allocs/op
BenchmarkAtomicWriteRead-8          61284550     23.73 ns/op    0 B/op    0 allocs/op
BenchmarkMutexReadOnly-8            16892307     71.56 ns/op    0 B/op    0 allocs/op
BenchmarkRWMutexReadOnly-8          28671775     51.60 ns/op    0 B/op    0 allocs/op
BenchmarkAtomicReadOnly-8           1000000000      0.11 ns/op    0 B/op    0 allocs/op
BenchmarkReadHeavyStoreRWMutex-8    14013139     81.74 ns/op    0 B/op    0 allocs/op
BenchmarkReadHeavyStoreMutex-8      14232532     86.07 ns/op    0 B/op    0 allocs/op
BenchmarkWriteHeavyStoreMutex-8      8637974    171.4 ns/op   18 B/op    1 allocs/op
BenchmarkWriteHeavyStoreRWMutex-8    8169729    181.1 ns/op   19 B/op    1 allocs/op
BenchmarkAtomicInt64Add-8           638972248      1.89 ns/op    0 B/op    0 allocs/op
BenchmarkAtomicInt64Load-8          1000000000      0.27 ns/op    0 B/op    0 allocs/op
BenchmarkCondSignal-8                3875018    310.6 ns/op   24 B/op    1 allocs/op
BenchmarkCondBroadcast-8                27349   44317 ns/op 2402 B/op  100 allocs/op
BenchmarkCondSignalChain-8           4463025    264.9 ns/op   24 B/op    1 allocs/op
```

关键发现：
- **RWMutex vs Mutex** - 读多写少时 RWMutex 快约 2x
- **Atomic vs Mutex** - 原子操作比互斥锁快约 5x（写入）到 300x（读取）
- **写多读少** - 使用普通 Mutex 性能更好
- **读多写少** - 使用 RWMutex 性能更好

## 实战建议

1. **优先使用原子操作** - 对于简单计数器和标志位，使用 atomic
2. **读多写少用 RWMutex** - 多线程读取，单线程写入时 RWMutex 性能更好
3. **写多读少用 Mutex** - 写入频繁时 RWMutex 额外开销不划算
4. **使用 sync.Once** - expensive initialization 使用单次执行保证

## 运行测试

```bash
go test -bench="Mutex|RWMutex|Atomic|Once" -benchmem ./code/sync/
```
