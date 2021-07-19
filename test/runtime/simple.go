package main

import (
	"fmt"
	"time"
)

func Producer(factor int,out chan <-int)  {
	for i := 0; ; i++ {
		out <- i*factor
	}
}

func Customer(ints <-chan int)  {
	for  i := range ints {
		fmt.Println(i)
	}
}

func main() {
	done := make(chan int, 10)

	go Producer(3,done)
	go Producer(5 , done)
	go Customer(done)

	time.Sleep(1*time.Millisecond)
}
