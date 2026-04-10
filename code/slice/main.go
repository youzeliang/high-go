package sliceopt

// PreAllocateSlice demonstrates the performance benefit of pre-allocating slices
func PreAllocateSlice(n int) []int {
	s := make([]int, 0, n)
	for i := 0; i < n; i++ {
		s = append(s, i)
	}
	return s
}

// DynamicSlice demonstrates dynamic slice growth without pre-allocation
func DynamicSlice(n int) []int {
	var s []int
	for i := 0; i < n; i++ {
		s = append(s, i)
	}
	return s
}

// SmallInitialSlice demonstrates starting with a small capacity
func SmallInitialSlice(n int) []int {
	s := make([]int, 0, 10)
	for i := 0; i < n; i++ {
		s = append(s, i)
	}
	return s
}

// CorrectCapacity demonstrates using the correct initial capacity
func CorrectCapacity(n int) []int {
	s := make([]int, n)
	for i := 0; i < n; i++ {
		s[i] = i
	}
	return s
}

// AppendMultiple demonstrates efficient append with multiple elements
func AppendMultiple(s []int, items []int) []int {
	return append(s, items...)
}

// CopySlice demonstrates efficient slice copying
func CopySlice(s []int) []int {
	result := make([]int, len(s))
	copy(result, s)
	return result
}

// FilterSlice demonstrates efficient in-place slice filtering
func FilterSlice(s []int, predicate func(int) bool) []int {
	result := s[:0]
	for _, v := range s {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// InefficientMapInLoop creates a map inside a loop for each iteration
func InefficientMapInLoop(n int) []map[string]int {
	result := make([]map[string]int, n)
	for i := 0; i < n; i++ {
		// Creating a new map each iteration - expensive!
		m := make(map[string]int)
		m["key"] = i
		result[i] = m
	}
	return result
}

// EfficientMapOutsideLoop creates map outside loop and reuses it
func EfficientMapOutsideLoop(n int) []map[string]int {
	result := make([]map[string]int, n)
	reusable := make(map[string]int)
	for i := 0; i < n; i++ {
		// Clone the reusable map
		m := make(map[string]int, len(reusable))
		for k, v := range reusable {
			m[k] = v
		}
		m["key"] = i
		result[i] = m
	}
	return result
}

// InefficientSliceCopy demonstrates inefficient repeated copying
func InefficientSliceCopy(src []int) []int {
	result := make([]int, 0, len(src))
	for _, v := range src {
		temp := make([]int, len(result)+1)
		copy(temp, result)
		temp[len(result)] = v
		result = temp
	}
	return result
}

// EfficientSliceCopy demonstrates efficient copying using built-in copy
func EfficientSliceCopy(src []int) []int {
	result := make([]int, len(src))
	copy(result, src)
	return result
}

// InefficientMapOperations demonstrates multiple individual map operations
func InefficientMapOperations(keys []string, values []int) map[string]int {
	m := make(map[string]int)
	for i := 0; i < len(keys); i++ {
		m[keys[i]] = values[i]
	}
	return m
}

// EfficientBatchMapOperation demonstrates batch map operation
func EfficientBatchMapOperation(keys []string, values []int) map[string]int {
	m := make(map[string]int, len(keys))
	for i := 0; i < len(keys); i++ {
		m[keys[i]] = values[i]
	}
	return m
}

// InefficientNestedSliceAccess demonstrates inefficient nested slice access
func InefficientNestedSliceAccess(matrix [][]int, row, col int) int {
	sum := 0
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			sum += matrix[i][j]
		}
	}
	return sum
}

// EfficientSliceFlatten demonstrates flattening nested slice for cache efficiency
func EfficientSliceFlatten(matrix [][]int, row, col int) int {
	// Flatten the matrix for better cache locality
	flat := make([]int, 0, row*col)
	for i := 0; i < row; i++ {
		flat = append(flat, matrix[i]...)
	}
	sum := 0
	for i := 0; i < len(flat); i++ {
		sum += flat[i]
	}
	return sum
}
