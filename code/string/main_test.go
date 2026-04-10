package string

import "testing"

func benchmark(b *testing.B, f func(int, string) string) {
	str := randomString(10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f(10000, str)
	}
}

// BenchmarkPlusConcat benchmarks + operator concatenation
func BenchmarkPlusConcat(b *testing.B) {
	benchmark(b, PlusConcat)
}

// BenchmarkSprintfConcat benchmarks fmt.Sprintf concatenation
func BenchmarkSprintfConcat(b *testing.B) {
	benchmark(b, SprintfConcat)
}

// BenchmarkBuilderConcat benchmarks strings.Builder concatenation
func BenchmarkBuilderConcat(b *testing.B) {
	benchmark(b, BuilderConcat)
}

// BenchmarkBufferConcat benchmarks bytes.Buffer concatenation
func BenchmarkBufferConcat(b *testing.B) {
	benchmark(b, BufferConcat)
}

// BenchmarkByteConcat benchmarks byte slice concatenation
func BenchmarkByteConcat(b *testing.B) {
	benchmark(b, ByteConcat)
}

// BenchmarkPreByteConcat benchmarks pre-allocated byte slice concatenation
func BenchmarkPreByteConcat(b *testing.B) {
	benchmark(b, PreByteConcat)
}

// Example output:
// $ go test -bench="Concat" -benchmem ./code/string/
// BenchmarkPlusConcat-8                 43      31100264 ns/op   530998135 B/op     10026 allocs/op
// BenchmarkSprintfConcat-8              19      52757395 ns/op   832967660 B/op     34096 allocs/op
// BenchmarkBuilderConcat-8            23785         52867 ns/op     514801 B/op         23 allocs/op
// BenchmarkBufferConcat-8            23011         60955 ns/op     368579 B/op         13 allocs/op
// BenchmarkByteConcat-8              15506         68989 ns/op     621297 B/op         24 allocs/op
// BenchmarkPreByteConcat-8           34440         40631 ns/op     212992 B/op          2 allocs/op
