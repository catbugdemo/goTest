package main

import (
	"fmt"
	"time"
)

// 生产者
func main() {
	inCh := make(chan int, 1)
	go func() {
		i := producer()
		inCh <-i
	}()

	select {
	case tmp:=<-inCh:
		fmt.Println(tmp)
		close(inCh)
	case <-time.After(time.Second*2):
		fmt.Println("运行超时")
	}
}

func producer() int {
	return 1
}