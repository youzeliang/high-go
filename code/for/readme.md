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


