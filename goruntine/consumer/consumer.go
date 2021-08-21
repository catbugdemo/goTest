package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

var wg sync.WaitGroup

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