package main

import (
	"fmt"
	"log"
	"time"
)

// 情形二：一个接收者和N个发送者，此唯一接收者通过关闭一个额外的信号通道来通知发送者不要再发送数据了。
// 1.
// n 个 producer 发送者
func producer(closed <-chan struct{}, nums ...int) <-chan int {
	inCh := make(chan int, len(nums))
	go func() {
		for {
			select {
			case <-closed:
				time.Sleep(time.Second * 2)
				close(inCh)
				log.Println("关闭通道")
				return
			default:
				for _, num := range nums {
					inCh <- num
				}
			}
		}
	}()
	return inCh
}

func main() {
	fmt.Println("hello world")
}
