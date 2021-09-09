package main

import (
	"fmt"
)

// 客户端
func main() {
	fmt.Println("开启客户端")
	//createSocket()
}

/*func createSocket()  {
	// 解析服务端地址
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("net.ResolveTCPAddr error:",err)
		return
	}
	// laddr 本地地址 , raddr 远程地址
	tcp, err := net.DialTCP("tcp", nil, addr)
	if err != nil {

	}

}*/
