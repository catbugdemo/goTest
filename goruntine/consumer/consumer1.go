package main

import (
	"fmt"
	"sync"
)

/*var wg sync.WaitGroup

func main()  {
	ints := make(chan int, 5)
	ints <-1
	ints <-2
	wg.Add(1)
	go consumer(ints)
	defer func() {
		close(ints)
		wg.Wait()
	}()
}


func consumer(ins chan int) {
	defer func() {
		wg.Done()
	}()
	for  {
		select {
		case i,ok:=<-ins:
			if !ok {
				return
			}
			// do something
			fmt.Println(i)
		case <-time.After(time.Second):
			log.Println("超时")
			return
		}
	}
}

// 出现该问题的原因是


// 实践
// all goroutines are asleep - deadlock!
*/

// M个接收者和一个发送者，发送者通过关闭用来传输数据的通道来传递发送结束信号


type mc struct {
	cond *sync.Cond
	done bool
}

func New() *mc {
	return &mc{
		cond: sync.NewCond(&sync.Mutex{}),
		done: false,
	}
}

// producer 一个发送者
func (m *mc) producer(nums ...int) <-chan int {
	inCh := make(chan int, len(nums))
	go func() {
		m.cond.L.Lock()
		defer func() {
			close(inCh)
			m.cond.L.Unlock()
			m.cond.Broadcast()
			m.done = true
		}()
		for _, num := range nums {
			inCh <- num
		}
	}()
	return inCh
}

func (m *mc) consumer(inCh <-chan int) <-chan int {
	outCh := make(chan int, len(inCh))
	go func() {
		m.cond.L.Lock()
		defer func() {
			defer m.cond.L.Unlock()
			close(outCh)
		}()

		for !m.done {
			m.cond.Wait()
		}
		for ch := range outCh {
			outCh <- ch
		}
	}()
	return inCh
}

func merge(chs ...<-chan int, ) <-chan int {
	var wg sync.WaitGroup
	// 将所有数据最后集中到一个通道中
	outCh := make(chan int, len(chs))

	// 将所有数据回收
	collect := func(in <-chan int) {
		defer wg.Done()
		for n := range in {
			// 将数据传入通道
			outCh <- n
		}
	}

	wg.Add(len(chs))
	for _, ch := range chs {
		go collect(ch)
	}
	go func() {
		wg.Wait()
		close(outCh)
	}()

	return outCh
}

func main() {
	m := New()
	inCh := m.producer(1, 2, 3, 4, 5, 6)
	out1 := m.consumer(inCh)
	out2 := m.consumer(inCh)
	out3 := m.consumer(inCh)

	for i := range merge(out1, out2, out3) {
		fmt.Println(i)
	}
}
