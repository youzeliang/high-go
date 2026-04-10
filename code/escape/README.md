# 栈分配 vs 堆分配 (Stack vs Heap Allocation)

## 概述

Go 编译器通过逃逸分析（Escape Analysis）自动决定变量应该分配在栈上还是堆上。栈分配比堆分配快得多，且不需要 GC 管理。

## 逃逸分析规则

### 会导致变量逃逸到堆的情况

1. **返回指针** - 函数返回局部变量的指针
2. **接口类型** - 变量赋值给 `interface{}` 或接口类型
3. **不确定的大小** - 编译器无法确定大小的slice/map
4. **闭包捕获** - 闭包引用外部变量
5. **发送指针到channel** - 将指针发送到channel

### 保持在栈上的情况

1. **确定大小的数组** - 固定大小的数组
2. **局部变量不逃逸** - 没有指针返回或接口赋值
3. **直接传递** - 通过值传递的结构体

## 性能对比

```
BenchmarkStackAlloc-8              1000000000    0.87 ns/op    0 B/op    0 allocs/op
BenchmarkHeapAlloc-8               147658928    8.17 ns/op    8 B/op    1 allocs/op
BenchmarkHeapAllocByInterface-8    695016416    1.64 ns/op    0 B/op    0 allocs/op
BenchmarkStackWithSlice-8           1483624    801.8 ns/op    0 B/op    0 allocs/op
BenchmarkHeapWithSlice-8            62760274   20.46 ns/op   80 B/op    1 allocs/op
BenchmarkStackCopy-8              1000000000    0.84 ns/op    0 B/op    0 allocs/op
BenchmarkPreAllocSlice-8            67722795   17.48 ns/op   80 B/op    1 allocs/op
BenchmarkPassByPointer-8          1000000000    0.81 ns/op    0 B/op    0 allocs/op
BenchmarkStackByValue-8           374701897    3.25 ns/op    0 B/op    0 allocs/op
BenchmarkSimpleAdd-8              1000000000    0.27 ns/op    0 B/op    0 allocs/op
BenchmarkInliningThreshold-8        36646768   32.97 ns/op    0 B/op    0 allocs/op
```

栈分配比堆分配快 **~10倍**！简单函数被内联后几乎零开销。

## 查看逃逸分析

使用 `-gcflags="-m"` 查看编译器的逃逸分析决策：

```bash
go build -gcflags="-m" ./code/escape/
```

输出示例：
```
./main.go:9:6: can inline StackAlloc
./main.go:14:6: can inline HeapAlloc
./main.go:15:9: &result escapes to heap
```

## 编译器优化

### -gcflags 优化等级

Go 编译器提供多个优化级别：

```bash
# 默认优化
go build ./...

# 显示所有优化决策
go build -gcflags="-m -m" ./...

# 禁用内联
go build -gcflags="-l" ./...

# 更激进禁用内联（用于测试）
go build -gcflags="-l -l" ./...

# 全部优化禁用（用于调试）
go build -gcflags="-N" ./...
```

### 内联优化

函数内联可以消除函数调用开销：

```go
// 简单函数会被内联
func SimpleAdd(a, b int) int {
    return a + b
}
```

**会被内联的条件：**
- 函数体较小
- 没有复杂控制流
- 不是递归函数

**不会被内联的函数：**
- 使用 `//go:noinline` 标记的函数
- 递归函数
- 函数体过大的函数

```go
//go:noinline
func NoInlineFunction() int {
    // 不会被内联
}
```

### 死码消除

编译器自动移除永远不会执行的代码：

```go
func DeadCodeElimination(flag bool) int {
    if flag {
        return 10
    }
    // 当 flag 已知为 true 时，这部分会被消除
    return DeadCodeElimination(true) + 20
}
```

## 性能优化技巧

### 减少指针解引用

```go
// 差：多次指针解引用
func ExpensivePointer(data *struct{ value int }) int {
    return data.value + data.value + data.value
}

// 好：缓存解引用后的值
func CachedPointer(data *struct{ value int }) int {
    v := data.value
    return v + v + v
}
```

### 局部变量缓存

```go
// 差：每次循环调用 len()
func WithoutCache(arr []int) int {
    sum := 0
    for i := 0; i < len(arr); i++ {
        sum += arr[i]
    }
    return sum
}

// 好：缓存 len()
func WithCache(arr []int) int {
    sum := 0
    ln := len(arr) // 缓存长度
    for i := 0; i < ln; i++ {
        sum += arr[i]
    }
    return sum
}
```

### 批量操作

```go
// 批量处理减少函数调用开销
func BatchOperation(items []int) int {
    sum := 0
    batchSize := 4
    for i := 0; i < len(items); i += batchSize {
        end := i + batchSize
        if end > len(items) {
            end = len(items)
        }
        for j := i; j < end; j++ {
            sum += items[j]
        }
    }
    return sum
}
```

## 优化建议

1. **避免不必要的指针** - 优先使用值传递
2. **固定大小优先** - 使用固定大小数组而非slice
3. **减少接口使用** - 在性能关键路径避免 `interface{}`
4. **预分配容量** - 为slice/map预分配容量减少扩容
5. **使用局部变量缓存** - 避免重复计算
6. **使用内联函数** - 简单函数会被自动内联

## 运行测试

```bash
# 逃逸分析
go test -bench="Alloc" -benchmem -run=^$ ./code/escape/

# 编译器优化
go test -bench="SimpleAdd|Inlining|Local" -benchmem -run=^$ ./code/escape/
```
