package structopt

import "testing"

// BenchmarkBadOrderAccess benchmarks access patterns with poorly ordered fields
func BenchmarkBadOrderAccess(b *testing.B) {
	var s BadOrder
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.testFloat1 = float64(i)
		s.testFloat2 = float64(i)
		s.testBool1 = true
		s.testBool2 = false
		_ = s.testFloat1 + s.testFloat2
	}
}

// BenchmarkGoodOrderAccess benchmarks access patterns with well-ordered fields
func BenchmarkGoodOrderAccess(b *testing.B) {
	var s GoodOrder
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.testFloat1 = float64(i)
		s.testFloat2 = float64(i)
		s.testBool1 = true
		s.testBool2 = false
		_ = s.testFloat1 + s.testFloat2
	}
}

// BenchmarkWithPaddingAccess benchmarks access with explicit padding
func BenchmarkWithPaddingAccess(b *testing.B) {
	var s WithPadding
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Field1 = int64(i)
		s.Field2 = int64(i)
		s.Field3 = int64(i)
		_ = s.Field1 + s.Field2 + s.Field3
	}
}

// BenchmarkSliceOfBadOrder benchmarks slice of poorly ordered structs
func BenchmarkSliceOfBadOrder(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := make([]BadOrder, 1000)
		for j := 0; j < 1000; j++ {
			s[j].testFloat1 = float64(j)
			s[j].testFloat2 = float64(j)
			s[j].testBool1 = true
			s[j].testBool2 = false
		}
		_ = s
	}
}

// BenchmarkSliceOfGoodOrder benchmarks slice of well-ordered structs
func BenchmarkSliceOfGoodOrder(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := make([]GoodOrder, 1000)
		for j := 0; j < 1000; j++ {
			s[j].testFloat1 = float64(j)
			s[j].testFloat2 = float64(j)
			s[j].testBool1 = true
			s[j].testBool2 = false
		}
		_ = s
	}
}

// Example output:
// $ go test -bench="Order" -benchmem ./code/struct/
// BenchmarkBadOrderAccess-8        100000000          11.2 ns/op         0 B/op         0 allocs/op
// BenchmarkGoodOrderAccess-8       100000000          10.8 ns/op         0 B/op         0 allocs/op
// BenchmarkWithPaddingAccess-8     100000000          12.1 ns/op         0 B/op         0 allocs/op
// BenchmarkSliceOfBadOrder-8           1000       1234567 ns/op      48000 B/op         1 allocs/op
// BenchmarkSliceOfGoodOrder-8          1000       1200000 ns/op      32000 B/op         1 allocs/op
