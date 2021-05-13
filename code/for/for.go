package _for

import (
	"math/rand"
	"time"
)

func generateWithCap(n int) []int {

	rand.Seed(time.Now().UnixNano())
	nums := make([]int, 0, n)
	for i := 0; i < n; i++ {
		nums = append(nums, rand.Int())
	}
	return nums

}
