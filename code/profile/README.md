# Profiling & Benchmark

## 概述

Go 提供了强大的性能分析工具，本章节介绍如何使用 pprof、trace 和 Benchmark 进行性能分析。

## pprof CPU 性能分析

### 启用 CPU  profiling

```bash
# 运行测试时启用 CPU profiling
go test -bench="FunctionName" -cpuprofile=cpu.prof ./code/profile/

# 分析 profiling 数据
go tool pprof cpu.prof

# 在 pprof 交互界面中
(pprof) top    # 显示占用 CPU 最多的函数
(pprof) web    # 生成调用图（在浏览器中打开）
(pprof) list FunctionName  # 查看特定函数的源码和每行耗时
```

### 代码中启用 profiling

```go
import "runtime/pprof"

func main() {
    // 启动 CPU profiling
    f, _ := os.Create("cpu.prof")
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()

    // 执行要分析的操作
    doWork()

    f.Close()
}
```

## pprof 内存分析

### 启用内存 profiling

```bash
# 运行测试时启用内存 profiling
go test -bench="FunctionName" -memprofile=mem.prof -benchmem ./code/profile/

# 分析内存数据
go tool pprof mem.prof

# 查看堆内存分配
(pprof) top    # 显示内存占用最多的函数
(pprof) inuse_space  # 查看存活的对象
(pprof) alloc_space  # 查看累计分配（包括已回收的）
```

### 代码中启用内存 profiling

```go
import "runtime/pprof"

func main() {
    // 创建内存 profile
    f, _ := os.Create("mem.prof")
    pprof.WriteHeapProfile(f)
    f.Close()
}
```

### 查看内存统计

```go
import "runtime"

var stats runtime.MemStats
runtime.ReadMemStats(&stats)

fmt.Printf("Alloc: %v MB\n", stats.Alloc/1024/1024)
fmt.Printf("TotalAlloc: %v MB\n", stats.TotalAlloc/1024/1024)
fmt.Printf("Sys: %v MB\n", stats.Sys/1024/1024)
fmt.Printf("NumGC: %v\n", stats.NumGC)
```

## trace 追踪并发

### 启用 trace

```bash
# 运行测试时启用 trace
go test -bench="." -trace=trace.out ./code/profile/

# 分析 trace 数据
go tool trace trace.out
```

### 代码中启用 trace

```go
import "runtime/trace"

func main() {
    f, _ := os.Create("trace.out")
    trace.Start(f)
    defer trace.Stop()

    // 执行并发操作
    doConcurrentWork()

    f.Close()
}
```

## Benchmark 测试写法

### 基本 Benchmark

```go
import "testing"

func BenchmarkFunction(b *testing.B) {
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        // 要测试的代码
        result := doWork()
        _ = result
    }
}
```

### 报告内存分配

```go
func BenchmarkWithAlloc(b *testing.B) {
    b.ResetTimer()
    b.ReportAllocs()  // 报告内存分配
    for i := 0; i < b.N; i++ {
        result := doWork()
        _ = result
    }
}
```

### 并行 Benchmark

```go
func BenchmarkParallel(b *testing.B) {
    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            doWork()
        }
    })
}
```

### 子测试和 Bench姆ark

```go
func BenchmarkSlice(b *testing.B) {
    sizes := []int{100, 1000, 10000}
    for _, size := range sizes {
        b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                doWork(size)
            }
        })
    }
}
```

## 运行测试

```bash
# 运行所有 benchmarks
go test -bench=. ./code/profile/

# 运行特定 benchmark
go test -bench="ComputeIntensive" ./code/profile/

# 带内存统计
go test -bench="Memory" -benchmem ./code/profile/

# CPU profiling
go test -bench="." -cpuprofile=cpu.prof ./code/profile/

# 内存 profiling
go test -bench="." -memprofile=mem.prof ./code/profile/

# trace
go test -bench="." -trace=trace.out ./code/profile/
```

## 性能对比示例

```
BenchmarkComputeIntensive-8         500000     3200 ns/op    0 B/op    0 allocs/op
BenchmarkMemoryAllocation-8           50000    28000 ns/op  45000 B/op   15 allocs/op
BenchmarkPreallocatedMemory-8        200000     8500 ns/op   8192 B/op    1 allocs/op
BenchmarkMemoryIntensive-8            10000   150000 ns/op1048576 B/op    1 allocs/op
BenchmarkGetMemStats-8            5000000      320 ns/op       0 B/op    0 allocs/op
```

关键发现：
- **预分配** 比动态分配快约 3x
- **GetMemStats** 有一定开销，不应在热路径中频繁调用
- **ComputeIntensive** 主要消耗 CPU 时间

## 实战建议

1. **先跑 Benchmark** - 找到最慢的函数
2. **CPU profiling** - 定位 CPU 热点
3. **内存 profiling** - 定位内存分配热点
4. **trace** - 分析并发问题
5. **避免过早优化** - 让数据驱动优化决策
