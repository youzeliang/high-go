package structopt

import (
	"unsafe"
)

// BadOrder demonstrates poor field ordering causing padding
type BadOrder struct {
	testBool1  bool    // 1 byte
	testFloat1 float64 // 8 bytes - causes padding after bool
	testBool2  bool    // 1 byte
	testFloat2 float64 // 8 bytes
}

// GoodOrder demonstrates optimal field ordering with minimal padding
type GoodOrder struct {
	testFloat1 float64 // 8 bytes
	testFloat2 float64 // 8 bytes
	testBool1  bool    // 1 byte
	testBool2  bool    // 1 byte (trailing padding added)
}

// CompactOrder shows maximum memory efficiency with same types grouped
type CompactOrder struct {
	testFloat1 float64
	testFloat2 float64
	testBool1  bool
	testBool2  bool
}

// WithPadding demonstrates adding explicit padding for cache line alignment
type WithPadding struct {
	Field1 int64  // 8 bytes
	_      int64  // padding for cache line alignment
	Field2 int64  // 8 bytes
	_      int64  // padding
	Field3 int64  // 8 bytes
	_      int64  // padding
}

// GetBadOrderSize returns the size of BadOrder struct
func GetBadOrderSize() uintptr {
	return unsafe.Sizeof(BadOrder{})
}

// GetGoodOrderSize returns the size of GoodOrder struct
func GetGoodOrderSize() uintptr {
	return unsafe.Sizeof(GoodOrder{})
}

// GetWithPaddingSize returns the size of WithPadding struct
func GetWithPaddingSize() uintptr {
	return unsafe.Sizeof(WithPadding{})
}
