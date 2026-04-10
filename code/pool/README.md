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