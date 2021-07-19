package teach

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestSimple(t *testing.T) {
	// 原因：开启了协程，但是主线程已经结束
	t.Run("testRun",
		func(t *testing.T) {


		go func() { // 初始化进程
			fmt.Println("Hello world")
		}()

			// time.Sleep(1 * time.Second)
		// 等待线程结束
		time.Sleep(1*time.Second)

	})
	// 原因：开启了协程，但是主线程已经结束


	// 优化 用 channel (目的，为了等待线程结束)
	t.Run("testLock", func(t *testing.T) {
		ints := make(chan int)  // 任何类型都行


		go func() {
			fmt.Println("Hello world")
			ints <- 1 // 添加
		}()

		<- ints  // 释放
	})


	// 优化，用 sync.WaitGroup
	t.Run("testWaitGroup", func(t *testing.T) {
		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("Hello world")
		}()

		wg.Wait()
	})
}
