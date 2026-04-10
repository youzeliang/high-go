# Slice 高效操作

## 概述

Go 的 slice 是动态数组，但不当的使用会导致频繁的内存分配和拷贝，影响性能。本章节展示如何高效地操作 slice，以及避免复杂的 slice/map 操作。

## 核心技巧

### 1. 预分配容量 (Pre-allocation)

使用 `make([]T, 0, capacity)` 预分配足够的容量，避免动态扩容。

```go
// 预分配
s := make([]int, 0, 10000)
for i := 0; i < 10000; i++ {
    s = append(s, i)
}

// 直接使用索引（最高效）
s := make([]int, 10000)
for i := 0; i < 10000; i++ {
    s[i] = i
}
```

### 2. 避免 nil slice append

```go
// 低效：nil slice 多次扩容
var s []int
for i := 0; i < n; i++ {
    s = append(s, i)
}

// 高效：预分配容量
s := make([]int, 0, n)
for i := 0; i < n; i++ {
    s = append(s, i)
}
```

### 3. 使用索引而非 append

```go
// 最高效：预分配 + 索引访问
s := make([]int, n)
for i := 0; i < n; i++ {
    s[i] = i
}
```

### 4. 原地过滤 (In-place filtering)

```go
// 高效：复用原 slice 底层数组
func FilterSlice(s []int, predicate func(int) bool) []int {
    result := s[:0]
    for _, v := range s {
        if predicate(v) {
            result = append(result, v)
        }
    }
    return result
}
```

### 5. 避免循环内创建 Map

```go
// 不推荐 - 每次循环创建新 map
for i := 0; i < n; i++ {
    m := make(map[string]int)
    m["key"] = i
}

// 推荐 - 减少 map 创建次数
reusable := make(map[string]int)
for i := 0; i < n; i++ {
    m := make(map[string]int, len(reusable))
    for k, v := range reusable {
        m[k] = v
    }
    m["key"] = i
}
```

### 6. 高效 Slice 复制

```go
// 不推荐 - 每次追加都重新分配和复制
result := make([]int, 0)
for _, v := range src {
    temp := make([]int, len(result)+1)
    copy(temp, result)
    temp[len(result)] = v
    result = temp
}

// 推荐 - 使用内置 copy
result := make([]int, len(src))
copy(result, src)
```

### 7. Map 预分配容量

```go
// 不推荐 - 多次扩容
m := make(map[string]int)
for i := 0; i < len(keys); i++ {
    m[keys[i]] = values[i]
}

// 推荐 - 预分配容量
m := make(map[string]int, len(keys))
for i := 0; i < len(keys); i++ {
    m[keys[i]] = values[i]
}
```

## 性能对比

```
BenchmarkPreAllocateSlice-8              177379     5742 ns/op   81920 B/op    1 allocs/op
BenchmarkDynamicSlice-8                  52528    21831 ns/op  357628 B/op   19 allocs/op
BenchmarkSmallInitialSlice-8             56910    21072 ns/op  350179 B/op   14 allocs/op
BenchmarkInefficientSliceCopy-8         314258     3818 ns/op   43120 B/op   101 allocs/op
BenchmarkEfficientSliceCopy-8          18898579       64.47 ns/op     896 B/op     1 allocs/op
BenchmarkInefficientMapOperations-8      25455    47662 ns/op  120334 B/op     30 allocs/op
BenchmarkEfficientBatchMapOperation-8    61299    19689 ns/op   57368 B/op      2 allocs/op
```

关键发现：
- **Slice 预分配** 比动态扩容快约 4x（5742ns vs 21831ns）
- **高效 slice copy** 比低效 copy 快约 60x（64ns vs 3818ns）
- **Map 批量操作** 比单独操作快约 2.4x（19689ns vs 47662ns）
- 正确预分配可减少 **~19x** 内存分配次数

## 运行测试

```bash
go test -bench="Slice|Map" -benchmem ./code/slice/
```

## 注意事项

1. **预估容量** - 如果知道大致大小，务必预分配
2. **容量翻倍** - slice 扩容时会创建新的底层数组，容量翻倍
3. **不要过度分配** - 过大的预分配会浪费内存
4. **slice 是引用类型** - 复制 slice 只是复制了指针和长度
5. **避免循环内创建 map** - 减少内存分配次数
