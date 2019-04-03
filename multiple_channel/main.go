package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func ToChansTimedContextGoroutines(ctx context.Context, d time.Duration, message int, ch ...chan int) (written int) {
	ctx, cancel := context.WithTimeout(ctx, d)
	defer cancel()
	var (
		wr int32
		wg sync.WaitGroup
	)
	wg.Add(len(ch))
	//这里的range ch，并不是从channel中读取数据，否则这块是会block的，因为没有数据可读
	//这里的range ch，实际上是从多个传进来的ch进行循环取到，不存在读取channel的操作
	for _, c := range ch {
		c := c
		//<-c
		go func() {
			defer wg.Done()
			select {
			case c <- message:
				atomic.AddInt32(&wr, 1)
			case <-ctx.Done():
			}
		}()
	}
	wg.Wait()
	return int(wr)
}

func main() {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch3 := make(chan int, 1)
	//如果这块需要往channel里面写数据，那就必须保证23行的读channel执行，否则会出现c <- message写不进去的问题
	//ch1 <- 1
	//ch2 <- 2
	//ch3 <- 3
	ret := ToChansTimedContextGoroutines(context.Background(), time.Second*3, 10, ch1, ch2, ch3)
	ret1 := <-ch1
	ret2 := <-ch2
	ret3 := <-ch3

	fmt.Println(ret)
	fmt.Println(ret1)
	fmt.Println(ret2)
	fmt.Println(ret3)

}
