# Ring Buffer 实现

## 概述

Ring Buffer（环形缓冲区）是一种固定大小的缓冲区，当数据达到末尾时会环绕到开头。非常适合 FIFO 队列场景，具有有界内存和高效的入队/出队操作。

## 核心模式

### 1. 基本 Ring Buffer

使用数组和原子计数器实现的线程安全的环形缓冲区。

```go
r := NewRingBuffer(1000)
r.Push(item)      // 添加元素
item, ok := r.Pop() // 取出元素
```

**特点：**
- 固定大小，有界内存
- 写入超过容量时自动覆盖最旧的数据
- 原子操作实现，无锁设计

### 2. 类型安全 Ring Buffer

针对特定类型优化的环形缓冲区，避免 interface{} 的装箱开销。

```go
r := NewIntRingBuffer(1000)
r.Push(42)
val, ok := r.Pop()
```

**特点：**
- 避免 interface{} 装箱
- 性能更好
- 类型安全

### 3. Channel-based Ring Buffer

使用 channel 实现的环形缓冲区，适合协程间通信。

```go
r := NewChannelRingBuffer(100)
ok := r.Push(item)  // 非阻塞写入
item, ok := r.Pop()  // 非阻塞读取
```

**特点：**
- 天然线程安全
- 支持非阻塞操作
- 适合生产者-消费者模式

## 性能对比

```
BenchmarkRingBufferPush-8           132110702     8.935 ns/op     8 B/op    0 allocs/op
BenchmarkRingBufferPop-8            100000000    10.02 ns/op     8 B/op    0 allocs/op
BenchmarkIntRingBufferPush-8        319798790     4.228 ns/op     0 B/op    0 allocs/op
BenchmarkIntRingBufferPop-8         272648698     5.154 ns/op     0 B/op    0 allocs/op
BenchmarkChannelRingBufferPush-8    128683969     9.290 ns/op     8 B/op    0 allocs/op
BenchmarkChannelRingBufferPop-8     44646109    26.76 ns/op     8 B/op    0 allocs/op
```

关键发现：
- **IntRingBuffer** 比普通 RingBuffer 快约 2x（避免 interface{} 开销）
- **RingBuffer Pop** 比 slice-based queue pop 快约 1000x（O(1) vs O(n)）
- **Channel Ring Buffer** 提供非阻塞操作，但 pop 较慢

## 实战建议

1. **使用 Ring Buffer** - 当你需要固定大小的 FIFO 队列时
2. **使用 IntRingBuffer** - 当处理数值类型时，获得更好性能
3. **使用 ChannelRingBuffer** - 当需要非阻塞操作和天然线程安全时
4. **避免 Slice Queue** - 出队操作 O(n)，性能差

## 运行测试

```bash
go test -bench="Ring" -benchmem ./code/ring/
```
