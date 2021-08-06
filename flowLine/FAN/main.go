package main

import (
	"fmt"
	"sync"
	"time"
)

// FAN-IN 和 FAN-OUT
// FAN模式可以让我们的流水线模型更好的利用 Golang并发
// 但是 FAN模式不一定是万能的，不见得能提高程序的性能
// 它不是万能的

// 任务分发是 FAN-OUT ,任务收集是 FAN-IN
// FAN-OUT 多个 goroutine 从同一个通道读取数据，直到该通道关闭
// FAN-IN  1个 goroutine 从多个通道读取数据，直到这些通道关闭

func producer(nums ...int) <-chan int {
	out := make(chan int, 10) // 不要使用无缓存通道,
	// 开启线程通道进行
	go func() {
		defer close(out)
		for _, num := range nums {
			out <- num
		}
	}()
	return out
}

func square(inCh <-chan int) <-chan int {
	out := make(chan int, 10)
	go func() {
		defer close(out)
		for n := range inCh {
			out <- n * n
		}
	}()
	return out
}

// 增加 merge() 入参是3个square各自写数据的通道，给这3个通道分别启动1个协程，把数据写入到自己创建的通道，并返回该通道，这是FAN-IN
// merge() 是进行数据分发 FAN-IN
func merge(cs ...<-chan int) <-chan int {
	out := make(chan int,10)

	var wg sync.WaitGroup

	collect := func(in <-chan int) {
		defer wg.Done()
		for n := range in {
			out <- n
		}
	}

	wg.Add(len(cs))
	// FAN-IN
	for _, c := range cs {
		go collect(c)
	}
	t := time.NewTimer(time.Microsecond*500)
	t.Stop()

	// 错误方式：直接等待是bug,死锁，因为merge写out,main却没有读
	// wg.Wait()
	// close(out)

	// 正确方式
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// 3 个 square 的协程是并发运行，结果顺序无法保证
// 修改 main()，启动3个square，这3个square 从 producer 生成的通道读数据，这是 FAN-OUT
func main() {
	in := producer(1, 2, 3, 4)
	c1 := square(in)
	c2 := square(in)
	c3 := square(in)

	// consumer
	for ret := range merge(c1, c2, c3) {
		fmt.Printf("%3d", ret)
	}
	fmt.Println()
}

// FAN 模式可以提高CPU利用率
// FAN 不一定能提升效率，降低程序运行时间