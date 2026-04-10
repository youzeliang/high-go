# Channel 缓冲与非缓冲

## 概述

Go 的 channel 是goroutine间通信的核心机制。理解缓冲与非缓冲channel的区别对于构建高效并发程序至关重要。

## 核心概念

### 非缓冲通道 (Unbuffered Channel)

```go
ch := make(chan int)  // 无缓冲区
```

**特性:**
- **同步模式** - 发送和接收必须同时准备
- **阻塞行为** - `ch <- value` 阻塞直到另一个goroutine接收
- **数据传递** - 发送者会阻塞直到接收者准备好接收费数据
- **内存** - 零缓冲区，实时handoff

**使用场景:**
- 需要严格同步的操作
- 确保消息准时送达
- 实现信号量或条件同步

### 缓冲通道 (Buffered Channel)

```go
ch := make(chan int, 100)  // 缓冲区大小100
```

**特性:**
- **异步模式** - 发送和接收可以独立进行
- **容量限制** - 缓冲区满时发送阻塞，缓冲区空时接收阻塞
- **性能提升** - 减少goroutine等待时间
- **解耦** - 生产者和消费者可以独立运行

**使用场景:**
- 生产者-消费者模式
- 批量处理
- 限流/背压控制

## 性能对比

```
BenchmarkUnbufferedSendReceive-8          5000000               281 ns/op
BenchmarkBufferedSendReceive-8            20000000                97.2 ns/op
```

**关键发现:**
- 缓冲通道发送/接收速度快约3倍
- 批量发送时差距更显著

### 批量发送性能

| 缓冲区大小 | 性能(ns/op) | 内存分配 |
|-----------|------------|----------|
| 1         | 1523       | 800 B/op |
| 10        | 272        | 8000 B/op |
| 100       | 146        | 80000 B/op |
| 1000      | 12.8       | 800000 B/op |

## select 语句

`select` 允许监听多个channel的非阻塞操作:

```go
select {
case v := <-ch1:
    return v, true
case v := <-ch2:
    return v, true
default:
    return 0, false  // 无数据可用，非阻塞
}
```

## Channel vs Mutex

| 场景 | Channel | Mutex |
|------|---------|-------|
| 数据传递 | ✓ 天然适合 | 需要额外设计 |
| 状态保护 | ✗ 不适合 | ✓ 适合 |
| 并发安全 | ✓ 内部安全 | ✓ 适合 |
| 性能 | 稍低(有开销) | 更高(直接操作) |

```go
// Mutex方式 - 适合状态保护
var counter int
var mu sync.Mutex
mu.Lock()
counter++
mu.Unlock()

// Channel方式 - 适合通信
ch := make(chan func(), 100)
ch <- func() { counter++ }
```

## 最佳实践

1. **选择合适的缓冲区大小**
   - 小流量: 1-10
   - 中流量: 100-1000
   - 大流量: 根据内存和延迟要求调整

2. **避免channel泄漏**
   - 确保发送和接收匹配
   - 使用context超时
   - 关闭不需要的channel

3. **流水线(Pipeline)设计**
   ```go
   // Stage 1: 读取
   // Stage 2: 处理
   // Stage 3: 输出
   input := Reader()
   processed := Processor(input)
   Writer(processed)
   ```

4. **Fan-Out模式**
   - 将任务分配给多个worker
   - 提高CPU利用率
   - 控制并发数量

## 运行测试

```bash
# 运行所有channel benchmarks
go test -bench="." -benchmem ./code/channel/

# 只测试缓冲vs非缓冲
go test -bench="Buffered" -benchmem ./code/channel/

# 测试Pipeline
go test -bench="Pipeline" -benchmem ./code/channel/

# 测试Mutex vs Channel
go test -bench="Counter" -benchmem ./code/channel/
```

## 注意事项

1. **缓冲区大小选择**: 过小无法发挥缓冲优势，过大浪费内存
2. **避免死锁**: 确保发送和接收操作最终都会执行
3. **关闭channel**: 只在不再发送时关闭，关闭已关闭的channel会panic
4. **nil channel**: 发送和接收nil channel会永久阻塞
