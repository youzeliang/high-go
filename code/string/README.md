# 字符串拼接优化

## 概述

Go 中字符串是不可变的，每次拼接都会创建新的字符串，导致大量内存分配。本章节展示各种字符串拼接方式的性能差异。

## 性能对比

| 方法 | ns/op | B/op | allocs/op | 推荐度 |
|------|-------|------|-----------|--------|
| PlusConcat (+) | 31100264 | 530998135 | 10026 | ❌ 不推荐 |
| SprintfConcat | 52757395 | 832967660 | 34096 | ❌ 严禁使用 |
| BuilderConcat | 52867 | 514801 | 23 | ✅ 推荐 |
| BufferConcat | 60955 | 368579 | 13 | ✅ 推荐 |
| ByteConcat | 68989 | 621297 | 24 | ⚠️ 一般 |
| PreByteConcat | 40631 | 212992 | 2 | ✅✅ 最佳 |

## 方法详解

### 1. + 运算符拼接 (不推荐)

```go
func PlusConcat(n int, str string) string {
    s := ""
    for i := 0; i < n; i++ {
        s += str
    }
    return s
}
```

**问题**: 每次拼接都创建新的字符串对象，产生 O(n²) 的时间复杂度。

### 2. fmt.Sprintf (严禁使用)

```go
func SprintfConcat(n int, str string) string {
    s := ""
    for i := 0; i < n; i++ {
        s = fmt.Sprintf("%s%s", s, str)
    }
    return s
}
```

**问题**: 性能最差，产生最多内存分配。

### 3. strings.Builder (推荐)

```go
func BuilderConcat(n int, str string) string {
    var builder strings.Builder
    for i := 0; i < n; i++ {
        builder.WriteString(str)
    }
    return builder.String()
}
```

**优点**: 内部使用 byte buffer，自动扩容。

### 4. bytes.Buffer (推荐)

```go
func BufferConcat(n int, s string) string {
    buf := new(bytes.Buffer)
    for i := 0; i < n; i++ {
        buf.WriteString(s)
    }
    return buf.String()
}
```

**优点**: 成熟稳定，适合大量字符串拼接。

### 5. 预分配 byte slice (最佳)

```go
func PreByteConcat(n int, str string) string {
    buf := make([]byte, 0, n*len(str))
    for i := 0; i < n; i++ {
        buf = append(buf, str...)
    }
    return string(buf)
}
```

**优点**: 只需 2 次内存分配，性能最佳。

## 性能提升

```
PlusConcat vs PreByteConcat:
- 速度提升: ~750x
- 内存节省: ~2500x
- 分配次数减少: ~5000x
```

## 运行测试

```bash
go test -bench="Concat" -benchmem ./code/string/
```

## 注意事项

1. **避免 + 拼接** - 在循环中绝对不要使用
2. **预分配容量** - 如果知道最终大小，使用 `Grow()` 或 `make([]byte, 0, capacity)`
3. **Builder vs Buffer** - strings.Builder 更轻量，但不适合少量拼接
4. **string() 转换** - 将 byte slice 转回 string 仍有成本