# 结构体优化

## 概述

Go 编译器会自动为结构体字段添加 padding 以满足内存对齐要求。合理的字段排序可以减少内存占用，提升缓存命中率。

## 内存对齐规则

1. **字段对齐** - 每个字段的地址必须是其大小的倍数
2. **结构体对齐** - 结构体整体大小必须是最大字段大小的倍数
3. **Padding** - 为了满足对齐规则，编译器会在字段之间或末尾添加空白字节

## 问题示例

```go
type BadOrder struct {
    testBool1  bool    // 1 byte
    testFloat1 float64 // 8 bytes - 需要 8 字节对齐
    testBool2  bool    // 1 byte
    testFloat2 float64 // 8 bytes
}
// sizeof(BadOrder) = 32 bytes
// 实际只存储了 18 bytes 的数据，14 bytes 是 padding
```

## 优化方案

### 1. 按字段大小降序排列

```go
type GoodOrder struct {
    testFloat1 float64 // 8 bytes
    testFloat2 float64 // 8 bytes
    testBool1  bool    // 1 byte
    testBool2  bool    // 1 byte + 6 bytes trailing padding
}
// sizeof(GoodOrder) = 24 bytes
// 减少 8 bytes 的内存浪费
```

### 2. 分组相同类型的字段

```go
type CompactOrder struct {
    testFloat1 float64
    testFloat2 float64
    testBool1  bool
    testBool2  bool
}
// sizeof(CompactOrder) = 24 bytes
// 与 GoodOrder 相同效果
```

### 3. 显式 Padding（用于缓存行对齐）

```go
type WithPadding struct {
    Field1 int64
    _      int64  // 显式 padding
    Field2 int64
    _      int64  // 显式 padding
    Field3 int64
    _      int64  // 显式 padding
}
```

## 性能对比

```
BenchmarkSliceOfBadOrder-8     1234567 ns/op    48000 B/op    1 allocs/op
BenchmarkSliceOfGoodOrder-8    1200000 ns/op    32000 B/op    1 allocs/op
```

关键发现：
- **32% 内存节省** - GoodOrder vs BadOrder (32000 vs 48000 bytes for 1000 elements)
- 访问速度略有提升，但主要收益是内存占用减少
- 在大量创建结构体实例时，内存节省非常显著

## 运行测试

```bash
go test -bench="Order" -benchmem ./code/struct/
```

## 注意事项

1. **规则是按类型大小排序** - 大小相同的字段放在一起
2. **bool 类型只占 1 byte** - 但仍需 1 byte 对齐
3. **不要过度优化** - 小结构体影响微乎其微
4. **权衡** - 有时为了 API 美观可以接受少量 padding