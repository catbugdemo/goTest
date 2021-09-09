package main

import "fmt"

// ## 协程的模板

// ## 前言
// 最近在写有关 go 的线程，但是突然一下子就卡住了
// 之前明明看过很多文章，也写过协程池，为什么一到实战就不行了呢，
// 好吧，其实是写的不多，但是快速上手也很重要，为什么不偷懒下直接套模板呢
// 以下是我总结的一些可以直接套用的简易模板

// ## 前置知识
// - go 协程
// - channel 通道
// - sync.WaitGroup

// 一个简单的开启线程

/*func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("hello")
	}()
	wg.Wait()
}*/

// producer 负责组装发送信息
func producer(nums ...int) <-chan int {
	inCh := make(chan int, len(nums))
	go func() {
		defer close(inCh)
		for _, num := range nums {
			inCh <- num
		}
	}()
	return inCh
}

// consumer 消费者
func consumer(inCh <-chan int) <-chan int {
	outCh := make(chan int, len(inCh))
	go func() {
		defer close(outCh)
		for in := range inCh {
			outCh <- in*in
		}
	}()
	return outCh
}

func main() {
	// 将数据组装为通道 -- 意味着可以被多组消费
	in := producer(1, 2, 3, 4)
	// 进行消费
	out := consumer(in)
	// 打印输出
	for i:= range out {
		fmt.Println(i)
	}
}
