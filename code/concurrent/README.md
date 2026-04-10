# Goroutine 高效调度

## 概述

Go 的并发模型以 goroutine 为基础，但不当使用会导致性能问题。本章节展示如何高效地使用 goroutine，包括 Worker Pool、信号量限制、原子操作等模式。

## 核心模式

### 1. Worker Pool (工作池)

Worker Pool 通过复用固定数量的 goroutine 来处理大量任务，避免为每个任务创建新的 goroutine。

```go
wp := NewWorkerPool(runtime.NumCPU(), 100)
wp.Start(func(t Task) Result {
    // 处理任务
    return Result{...}
})
for i := 0; i < 1000; i++ {
    wp.Submit(Task{ID: i, Data: i})
}
wp.Close()
wp.Wait()
```

**优点：**
- 限制并发 goroutine 数量
- 减少 goroutine 创建/销毁开销
- 提高 CPU 利用率

### 2. Semaphore (信号量)

使用信号量限制同时执行的 goroutine 数量，防止资源耗尽。

```go
sem := NewSemaphore(runtime.NumCPU())
var wg sync.WaitGroup
for i := 0; i < 100; i++ {
    wg.Add(1)
    go func(id int) {
        sem.Execute(func() {
            defer wg.Done()
            // 受保护的工作
        })
    }(i)
}
wg.Wait()
```

### 3. Atomic Counter (原子计数器)

在高并发场景下，使用原子操作比互斥锁更高效。

```go
counter := NewAtomicCounter()
// 并发 increment，无锁
counter.Increment()
counter.Add(10)
```

**性能对比：**
```
BenchmarkAtomicCounterIncrement-8    625M ops    1.9 ns/op    0 allocs/op
BenchmarkMutexCounterIncrement-8       5M ops  380 ns/op    0 allocs/op
```
- **原子操作比互斥锁快 ~200 倍**

### 4. Once 初始化

`sync.Once` 确保expensive initialization 只执行一次。

```go
init := NewExpensiveInit()
init.Get() // 只执行一次初始化
init.Get() // 直接返回缓存值
```

### 5. Context 超时控制

使用 context 实现 goroutine 的超时控制和取消。

```go
task := NewLongRunningTask(200 * time.Millisecond)
if task.Execute(func(ctx context.Context) error {
    // 执行工作，可检查 ctx.Done()
    return nil
}) {
    // 成功完成
} else {
    // 超时或被取消
}
```

## 性能对比

```
BenchmarkAtomicCounterIncrement-8      625M        1.9 ns/op      0 B/op    0 allocs/op
BenchmarkMutexCounterIncrement-8          5M      380 ns/op      0 B/op    0 allocs/op
BenchmarkSharedOnceInitialization-8    1000M        0.6 ns/op      0 B/op    0 allocs/op
```

关键发现：
- **原子操作** vs 互斥锁：约 **200x** 性能提升
- **Worker Pool** 减少 goroutine 创建开销
- **sync.Once** 保证单次初始化，多线程下高效

## 实战建议

1. **使用 Worker Pool** - 处理大量短任务时复用 goroutine
2. **限制并发数** - 使用信号量防止资源耗尽
3. **优先原子操作** - 简单计数等使用 atomic 而非 mutex
4. **Context 取消** - 长时间运行的任务使用 context 超时控制
5. **GOMAXPROCS** - 了解 CPU 核心数，合理配置

## 运行测试

```bash
go test -bench="Goroutine|Semaphore|Atomic|Once" -benchmem ./code/concurrent/
```