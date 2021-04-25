package main

//type TreeNode struct { 
//	Val int 
//	Left, Right *TreeNode 
//} 

//      5

type TreeNode struct {
	Val         int
	left, right *TreeNode
}

func main() {



	//type t struct {
	//}
	//
	//var i interface{}
	//
	//fmt.Println(i == nil)
	//i = (*t)(nil)
	//fmt.Println(i == nil)


	//var a = []int{1,2,3}
	//var b = []int{4,5,6}
	//
	//fmt.Println(KthElement(a,b,2))

	a := make([]int,0,5)

	a= append(a )

}

func KthElement(a []int, b []int, k int) int {

	
	left, right, t := len(a)-1, len(b)-1, len(a)+len(b)
	
	if k > left + right {
		return -1
	}

	c := make([]int,t)
	for left >= 0 && right >= 0 {

		if a[left] > b[right] {
			t--
			c[t] = a[left]
			left--
		} else {
			t--
			c[t] = b[right]
			right--
		}
	}

	for left >= 0 {
		t--
		c[t] = a[left]
		left--
	}

	for right >= 0 {
		t--
		c[t] = b[right]
		right--
	}

	return c[k-1]
}
