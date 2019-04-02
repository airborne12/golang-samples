//经典csp实现打印1-10000素数打印
package main

import (
	"fmt"
)

func main() {
	origin, wait := make(chan int), make(chan struct{})
	Processor(origin, wait)
	for num := 2; num < 10; num++ {
		origin <- num
	}
	close(origin)
	<-wait
}

//Processor 逻辑
/*
	输入： 2，3，4，5，6，7，8，9，10到channel origin
	goroutine 1 素数：2 过滤后输出到channel out：3，5，7，9
	goroutine 2 素数：3 过滤后输出到channel out: 5, 7
	goroutine 3 素数：5 过滤后输出到channel out：7
	goroutine 4 素数：7 过滤后关闭channel out
*/
func Processor(seq chan int, wait chan struct{}) {
	go func() {
		prime, ok := <-seq
		if !ok {
			close(wait)
			return
		}
		fmt.Println(prime)
		out := make(chan int)
		Processor(out, wait)
		for num := range seq {
			if num%prime != 0 {
				out <- num
			}
		}
		close(out)
	}()
}
