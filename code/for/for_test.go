package _for

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {

	ch := make(chan string)
	go func() {
		ch <- "Go"
		ch <- "语言"
		ch <- "高性能"
		ch <- "编程"
		close(ch)
	}()
	for n := range ch {
		fmt.Println(n)
	}
}