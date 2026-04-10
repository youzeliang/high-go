# sync.Pool 对象池

## 概述

`sync.Pool` 是 Go 标准库提供的临时对象池，用于缓存临时分配的对象以减少 GC 压力。它是在高并发场景下减少内存分配开销的重要工具。

## 核心特性

1. **线程安全** - sync.Pool 的 Get/Put 操作是协程安全的
2. **临时存储** - 对象可能在任何时候被自动移除（GC 时）
3. **无保证** - 无法保证对象会被返回（对象可能已被清除）
4. **懒加载** - New 函数在池为空时才被调用

## 使用场景

- 频繁分配和释放相同类型的对象（如 bytes.Buffer、[]byte）
- 减少 GC 压力，特别是在高吞吐量服务中
- 复用临时对象，如请求处理中的缓冲区

## 性能对比

```
BenchmarkNoPool-8             500000              3284 ns/op            4096 B/op        1 allocs/op
BenchmarkWithPool-8           5000000               286 ns/op              0 B/op        0 allocs/op
```

使用 sync.Pool 可以显著减少：
- 内存分配次数（0 allocs vs 1 allocs）
- 执行时间（~10x 提升）
- GC 压力

## 实现示例

### 基础字节池

```go
type Pool struct {
    slice sync.Pool
}

func NewPool() *Pool {
    return &Pool{
        slice: sync.Pool{
            New: func() interface{} {
                return make([]byte, 1024)
            },
        },
    }
}
```

### 通用对象池

```go
type ObjectPool struct {
    pool sync.Pool
}

func NewObjectPool(factory func() interface{}) *ObjectPool {
    return &ObjectPool{
        pool: sync.Pool{New: factory},
    }
}
```

## 注意事项

1. **不要 Put nil 值** - 放入池中的对象不应为 nil
2. **对象状态** - 取出的对象应该重置状态（如 bytes.Buffer 需要 Reset）
3. **不适合长期持有** - sync.Pool 中的对象随时可能被清除
4. **Pool.New 必须线程安全** - 避免竞态条件

## 运行测试

```bash
go test -bench="Pool$" -benchmem ./code/pool/
```

---

## 减少 GC 压力的高级池化技术

### 1. JSON 编码池 (JSONBufferPool)

在高吞吐量 JSON 处理场景中，复用 bytes.Buffer 可显著减少 GC 压力：

```go
type JSONBufferPool struct {
    pool sync.Pool
}

func NewJSONBufferPool() *JSONBufferPool {
    return &JSONBufferPool{
        pool: sync.Pool{
            New: func() interface{} { return &bytes.Buffer{} },
        },
    }
}
```

**性能对比:**
```
BenchmarkJSONEncodingNoPool-8      4256054  283.7 ns/op  352 B/op  4 allocs/op
BenchmarkJSONEncodingWithPool-8    5040988  238.5 ns/op  192 B/op  2 allocs/op
```
- 减少 50% 内存分配
- 提升约 20% 性能

### 2. 多尺寸对象池 (MultiSizeObjectPool)

针对不同大小的缓冲区维护专用池，减少内存碎片：

```go
type MultiSizeObjectPool struct {
    pools []sync.Pool
    sizes []int
}
```

**性能对比:**
```
BenchmarkMultiSizeNoPool-8    1934446  626.8 ns/op  1816 B/op  1 allocs/op
BenchmarkMultiSizePoolGet-8  2171092  576.8 ns/op    24 B/op  1 allocs/op
```

### 3. 行缓冲区池 (RowBufferPool)

数据库行处理场景，减少结构体分配：

```go
type RowBuffer struct {
    Columns []string
    Values  []interface{}
}
```

**性能对比:**
```
BenchmarkRowBufferNoPool-8    11245694   115.0 ns/op  512 B/op  2 allocs/op
BenchmarkRowBufferWithPool-8 129770016     8.4 ns/op    0 B/op  0 allocs/op
```
- **13x 性能提升**
- 零内存分配

## 实战建议

1. **热点对象优先** - 对于频繁分配/释放的对象考虑池化
2. **注意池大小** - 过大占用内存，过小效果不明显
3. **避免池中对象累积** - 及时回收不使用的对象
4. **结合实际场景** - JSON 处理、数据库行缓冲、临时缓冲区等

```bash
# 运行所有池化技术 benchmarks
go test -bench="." -benchmem ./code/pool/
```