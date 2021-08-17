package main

import "fmt"

// Golang 并发的核心思路是关注数据流动，数据流动的过程交给 channel，
// 数据处理的每个环节交给 gorountine，将这些流程画起来，有始有终形成一条线，那就能构成流水线模型

// 一个简单的例子：计算整数切片中元素的平方值并把它打印出来
//		非并发：遍历整个切片，然后计算平方，打印结果
//		并发：1.遍历切片，生产者。2.计算平方值。3.打印结果，消费者。

//	producer() 负责生产数据，将数据写入通道，并把它写数据的通道返回
//	square() 负责操作，负责从某个通道读数字，然后计算平方，将结果写入通道，并把它的输出通道返回
//	main() 负责启动 producer 和 square，同时也是消费者，读取suqre的结果，并打印出来

// producer 将数组切片存入通道中
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
			out <- n*n
		}
	}()
	return out
}

func main() {
	in := producer(1, 2, 3, 4)
	ch := square(in)

	// consumer
	for ret := range ch {
		fmt.Printf("%3d",ret)
	}
	fmt.Println()
}

// 流水线特点
//	1.每个阶段把数据通过channel传递给下一个阶段
//	2.每个阶段要创建1个goruntine和一个通道，这个gorountine向里面写数据，函数要返回这个通道
//	3.有一个函数来组织流水线，例子中时main函数