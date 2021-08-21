package main

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
