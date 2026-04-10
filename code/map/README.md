# Map 并发安全

## 概述

Go 的原生 map 不是线程安全的。在并发场景下，需要使用适当的同步机制来保护 map 的访问。本章节展示不同并发安全的 map 实现方式及其性能对比。

## 核心模式

### 1. Mutex Map

使用互斥锁保护 map 访问，简单但所有操作都需要获取锁。

```go
m := NewMutexMap()
m.Set("key", 1)
val, ok := m.Get("key")
```

**特点：**
- 实现简单
- 所有操作都需要获取锁
- 适用于读少写多的场景

### 2. RWMutex Map

使用读写互斥锁，允许多个读锁并发，写锁独占。

```go
m := NewRWMutexMap()
// 读操作 - 使用 RLock
val, ok := m.Get("key")
// 写操作 - 使用 Lock
m.Set("key", 1)
```

**特点：**
- 读操作可以并行
- 适用于读多写少的场景
- 性能优于普通 Mutex

### 3. sync.Map

Go 标准库提供的并发安全的 map，针对特定访问模式优化。

```go
m := NewSyncMap()
m.Set("key", 1)
val, ok := m.Get("key")
```

**特点：**
- 使用无需锁的读写分离设计
- 适合读多写少、只会增加 key 的场景
- 单一 key 频繁访问时性能较差

### 4. 批量操作

将多个操作合并为一次锁获取，减少锁竞争。

```go
m := NewBatchMutexMap()
// 批量写入 - 只获取一次锁
pairs := map[string]int{"a": 1, "b": 2, "c": 3}
m.SetBatch(pairs)
```

## 性能对比

```
BenchmarkMutexMapSingleKey-8         100000000    10.42 ns/op    0 B/op    0 allocs/op
BenchmarkRWMutexMapSingleKey-8       100000000    10.68 ns/op    0 B/op    0 allocs/op
BenchmarkSyncMapSingleKey-8           16344885    77.42 ns/op   40 B/op    2 allocs/op
BenchmarkMutexMapUniqueKeys-8         49677721    22.91 ns/op    6 B/op    1 allocs/op
BenchmarkSyncMapUniqueKeys-8          7617940   180.1 ns/op   67 B/op    4 allocs/op
BenchmarkRWMutexMapReadHeavy-8         800782   1563 ns/op    0 B/op    0 allocs/op
BenchmarkBatchMutexMapSet-8            686047   1719 ns/op    0 B/op    0 allocs/op
BenchmarkNaiveMutexMapSet-8            465049   2597 ns/op  400 B/op  100 allocs/op
BenchmarkSyncMapLoadFineGrained-8      572053   2136 ns/op    0 B/op    0 allocs/op
```

关键发现：
- **sync.Map** 单 key 访问比 mutex 慢约 7x
- **sync.Map** 适合 key 只会增加、读多写少的场景
- **RWMutex** 在读多写少时比普通 Mutex 性能更好
- **批量操作** 比多次单独操作效率更高

## 实战建议

1. **优先使用 sync.Map** - 当 key 只会增加、读远多于写时
2. **使用 RWMutex** - 读多写少，且需要频繁读取时
3. **使用 Mutex** - 写操作较多，或需要精确控制时
4. **批量操作** - 减少锁获取次数，提高性能

## 运行测试

```bash
go test -bench="Map" -benchmem ./code/map/
```
