package esc

import "fmt"

// goroutine 退出，三种优雅退出 goroutine 的方法
// 只要采用最佳实践去设计，基本上就可以确保 goroutine 退出上不会有问题，尽情享用

// 1. for-range
func F(in <-chan int) {
	go func(in <-chan int) {
		for x := range in {
			fmt.Printf("Process %d\n", x)
		}
	}(in)
}

// 2.使用 ,ok退出
// select 不会再 nil 的通道上进行等待
// 把只读通道设置为 nil 即可
func O(in <-chan int) {
	go func() {
		for {
			select {
			case x, ok := <-in:
				if !ok {
					return
				}
				fmt.Printf("Process %d\n",x)

			}
		}
	}()
}

// 3.使用 ,ok 来退出使用 for-select 协程
// 使用一个专门的通道，发送退出的信号，可以解决这类问题

// select 的3个进阶特性
// 1.nil的通道永远阻塞
// 2.如何跳出 for-select
// 3.select{} 永远阻塞阻塞，等价于
//		ch := make(chan int)
//		<-ch