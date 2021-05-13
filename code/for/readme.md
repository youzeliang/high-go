### for 和 range 的性能比较


#### 1. []int

```go

func BenchmarkName(b *testing.B) {

	nums := generateWithCap(1024 * 1024)

	for i := 0; i < b.N; i++ {
		l := len(nums)
		var tmp int
		for k := 0; k < l; k++ {
			tmp = nums[k]
		}
		_ = tmp
	}
}


func BenchmarkRangeIntSlice(b *testing.B) {
	nums := generateWithCap(1024 * 1024)
	for i := 0; i < b.N; i++ {
		var tmp int
		for _, num := range nums {
			tmp = num
		}
		_ = tmp
	}
}

```


运行结果

```shell
go test -bench=IntSlice$ .

cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkName-12                    2103            556512 ns/op
BenchmarkRangeIntSlice-12           2067            553682 ns/op


```

从结果可以看到，遍历 []int 类型的切片，for 与 range 性能几乎没有区别。


#### 2.  []struct

```go
func BenchmarkForStruct(b *testing.B) {
	var items [1024]Item
	for i := 0; i < b.N; i++ {
		length := len(items)
		var tmp int
		for k := 0; k < length; k++ {
			tmp = items[k].id
		}
		_ = tmp
	}
}

func BenchmarkRangeIndexStruct(b *testing.B) {
	var items [1024]Item
	for i := 0; i < b.N; i++ {
		var tmp int
		for k := range items {
			tmp = items[k].id
		}
		_ = tmp
	}
}

func BenchmarkRangeStruct(b *testing.B) {
	var items [1024]Item
	for i := 0; i < b.N; i++ {
		var tmp int
		for _, item := range items {
			tmp = item.id
		}
		_ = tmp
	}
}
```

```shell
go test -bench=Struct$ .

cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkForStruct-12                    4479685               271.0 ns/op
BenchmarkRangeIndexStruct-12             4505242               264.6 ns/op
BenchmarkRangeStruct-12                     6486            235165 ns/op
```

- 仅遍历下标的情况下，for 和 range 的性能几乎是一样的。
- items 的每一个元素的类型是一个结构体类型 Item，Item 由两个字段构成，一个类型是 int，一个是类型是 [4096]byte，也就是说每个 Item 实例需要申请约 4KB 的内存。
- 在这个例子中，for 的性能大约是 range (同时遍历下标和值) 的 2000 倍。


#### []int 和 []struct{} 的性能差异

与 for 不同的是，range 对每个迭代值都创建了一个拷贝。因此如果每次迭代的值内存占用很小的情况下，for 和 range 的性能几乎没有差异，但是如果每个迭代值内存占用很大，例如上面的例子中，每个结构体需要占据 4KB 的内存，这种情况下差距就非常明显了。

我们可以用一个非常简单的例子来证明 range 迭代时，返回的是拷贝。

#### []*struct{}

```go

func generateItems(n int) []*Item {
	items := make([]*Item, 0, n)
	for i := 0; i < n; i++ {
		items = append(items, &Item{id: i})
	}
	return items
}

func BenchmarkForPointer(b *testing.B) {
	items := generateItems(1024)
	for i := 0; i < b.N; i++ {
		length := len(items)
		var tmp int
		for k := 0; k < length; k++ {
			tmp = items[k].id
		}
		_ = tmp
	}
}

func BenchmarkRangePointer(b *testing.B) {
	items := generateItems(1024)
	for i := 0; i < b.N; i++ {
		var tmp int
		for _, item := range items {
			tmp = item.id
		}
		_ = tmp
	}
}

```


```shell

cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkForPointer-12            747474              1441 ns/op
BenchmarkRangePointer-12          777216              1342 ns/op

```

切片元素从结构体 Item 替换为指针 *Item 后，for 和 range 的性能几乎是一样的。


range 在迭代过程中返回的是迭代值的拷贝，如果每次迭代的元素的内存占用很低，那么 for 和 range 的性能几乎是一样，例如 []int。但是如果迭代的元素内存占用较高，例如一个包含很多属性的 struct 结构体，那么 for 的性能将显著地高于 range，有时候甚至会有上千倍的性能差异。对于这种场景，建议使用 for，如果使用 range，建议只迭代下标，通过下标访问迭代值，这种使用方式和 for 就没有区别了。如果想使用 range 同时迭代下标和值，则需要将切片/数组的元素改为指针，才能不影响性能。