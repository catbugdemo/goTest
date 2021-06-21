package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

//阻塞IO模型是使用 一个线程处理一个连接
func ListenAndServe(address string) {
	//返回一个在本地网络地址Iaddr上监听的Listener
	listener, e := net.Listen("tcp", address)
	if e != nil {
		log.Fatal(fmt.Sprintf("listen err:%v", e))
	}
	//注意关闭
	defer listener.Close()
	log.Println(fmt.Sprintf("bind:%s,start listening...", address))

	for {
		//Accept 会一直阻塞指导有新的连接建立或者listen终端才会返回
		conn, e := listener.Accept()
		if e != nil {
			// 通常是由于listener被关闭无法继续监听导致的错误
			log.Fatal(fmt.Sprintf("accept err:%v", e))
		}
		go Handler(conn)
	}
}

func Handler(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		// ReadString 会一直阻塞指导遇到分隔符 '\n'
		//遇到分隔符后 ReadString 会返回上次遇到分隔符到现在收到的所有数据
		//若在遇到分隔符之前发生异常 会返回已收到的数据和错误信息
		msg, e := reader.ReadString('\n')
		if e != nil {
			// 通常遇到的错误是连接中断或被关闭，用 io.EOF 表示
			if e == io.EOF {
				log.Println("connection close")
			}else {
				log.Println(e)
			}
			return
		}
		b := []byte(msg)
		// 将收到的信息发送给客户端
		conn.Write(b)
	}
}

func main() {
	ListenAndServe(":8000")
}
