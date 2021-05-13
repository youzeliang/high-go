package _for

import (
	"reflect"
	"testing"
)

type Item struct {
	id  int
	val [4096]byte
}

func Test_generateWithCap(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateWithCap(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateWithCap() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func BenchmarkName(b *testing.B) {
//
//	nums := generateWithCap(1024 * 1024)
//
//	for i := 0; i < b.N; i++ {
//		l := len(nums)
//		var tmp int
//		for k := 0; k < l; k++ {
//			tmp = nums[k]
//		}
//		_ = tmp
//	}
//}
//
//func BenchmarkRangeIntSlice(b *testing.B) {
//	nums := generateWithCap(1024 * 1024)
//	for i := 0; i < b.N; i++ {
//		var tmp int
//		for _, num := range nums {
//			tmp = num
//		}
//		_ = tmp
//	}
//}

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
