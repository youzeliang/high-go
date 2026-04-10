package escape

// StackAlloc demonstrates stack allocation (no escape to heap)
//go:noinline
func StackAlloc() int {
	x := 10
	y := 20
	return x + y
}

// HeapAlloc demonstrates heap allocation (escape to heap)
// The returned pointer escapes to the heap
//go:noinline
func HeapAlloc() *int {
	x := 10
	y := 20
	result := x + y
	return &result
}

// HeapAllocByInterface demonstrates how interface causes heap allocation
//go:noinline
func HeapAllocByInterface() interface{} {
	x := 10
	return x
}

// StackWithSlice demonstrates slice that stays on stack
//go:noinline
func StackWithSlice() [1024]int {
	var arr [1024]int
	for i := 0; i < 1024; i++ {
		arr[i] = i
	}
	return arr
}

// HeapWithSlice demonstrates slice header escaping to heap
//go:noinline
func HeapWithSlice() []int {
	// slice header escapes, but underlying array may not
	arr := make([]int, 10)
	for i := 0; i < 10; i++ {
		arr[i] = i
	}
	return arr
}

// StackCopy demonstrates efficient stack-to-stack copying
//go:noinline
func StackCopy() int {
	x := 10
	y := x
	return y
}

// PreAllocSlice demonstrates pre-allocated slice (stays on stack if small)
//go:noinline
func PreAllocSlice(n int) []int {
	// When n is known constant and small, compiler may optimize to stack
	s := make([]int, 0, 10)
	for i := 0; i < n && i < 10; i++ {
		s = append(s, i)
	}
	return s
}

// PassByPointer demonstrates passing pointer (may escape or not)
//go:noinline
func PassByPointer(p *int) int {
	return *p
}

// StackByValue demonstrates struct passed by value (copied on stack)
//go:noinline
func StackByValue(s [4]int) int {
	sum := 0
	for _, v := range s {
		sum += v
	}
	return sum
}

// =============================================================================
// Compiler Optimization Examples
// =============================================================================

// SimpleAdd is a simple function that the compiler can inline
func SimpleAdd(a, b int) int {
	return a + b
}

// InliningThreshold demonstrates a function right at the inlining threshold
// Functions with many statements or complex control flow may not be inlined
func InliningThreshold(n int) int {
	sum := 0
	for i := 0; i < n; i++ {
		sum += i
	}
	return sum
}

// NoInlineFunction uses go:noinline to prevent inlining
//go:noinline
func NoInlineFunction(n int) int {
	sum := 0
	for i := 0; i < n; i++ {
		sum += i
	}
	return sum
}

// DeadCodeElimination demonstrates unused code that the compiler removes
func DeadCodeElimination(flag bool) int {
	if flag {
		return 10
	}
	// This branch is never taken when flag is true
	// Compiler can eliminate the unreachable code
	return DeadCodeElimination(true) + 20
}

// UnreachableCode is never called - will be eliminated
func UnreachableCode() int {
	return 999
}

// CalledFunction is the only function actually used
func CalledFunction() int {
	return UnreachableCode() // This call might be optimized away too
}

// LocalVariableCaching demonstrates caching local variables
func LocalVariableCaching(arr []int) int {
	sum := 0
	// Cache len to avoid repeated calls
	ln := len(arr)
	for i := 0; i < ln; i++ {
		sum += arr[i]
	}
	return sum
}

// PointerChasing demonstrates multiple pointer dereferences (expensive)
func PointerChasing(data *struct {
	next *struct {
		value int
	}
}) int {
	// Multiple pointer dereferences - cache if possible
	return data.next.value
}

// BatchOperation demonstrates processing multiple items together
func BatchOperation(items []int) int {
	sum := 0
	// Process in batches for better cache utilization
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

// ReducePointerDerefs caches frequently accessed values
func ReducePointerDerefs(data *struct{ value int }) int {
	// Cache the dereferenced value in a local variable
	v := data.value
	sum := v + v + v // Multiple uses of cached value
	return sum
}