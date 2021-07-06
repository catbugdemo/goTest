package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

//使用一个计数器(应该是？)
var group = sync.WaitGroup{}

func main() {
	//当需要开启一个线程时，设定等待的线程+1
	group.Add(1)
	go UserInput()

	//等待全部线程结束
	group.Wait()
	os.Exit(0)
}

func UserInput() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("输入指令：")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		switch text {
		case "exit":
			{ // 退出本程序
				fmt.Print("程序退出\n")
				//线程数-1
				group.Done()
				return
			}
		default:
			fmt.Println("不可识别的指令")
		}
	}
}
