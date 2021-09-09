package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
)

var connSlice []*net.TCPConn

// 服务器
func main() {
	fmt.Println("服务端")

	// 建立长连接
	createTcp()
}

// 创建 tcp 长连接
func createTcp() {
	// 解析 tcp 服务
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("net.ResolveTCPAddr error:", err)
		return
	}
	// 监听指定服务
	tcp, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Println("net.ListenTCP error:", err)
		return
	}
	defer tcp.Close()
	for {
		// 阻塞，当有客户端连接时，才运行下面
		acceptTCP, err := tcp.AcceptTCP()
		if err != nil {
			fmt.Println("tcpListener error :", err)
			continue
		}
		fmt.Println("A client connected:", acceptTCP.RemoteAddr().String())
		boradcastMessage(acceptTCP.RemoteAddr().String() + "进入房间" + "\n")
		connSlice = append(connSlice, acceptTCP)
		// 开启一个协程处理信息
		go tcpPipe(acceptTCP)
	}
}

// 对客户端做出反应
func tcpPipe(conn *net.TCPConn) {
	ipStr := conn.RemoteAddr().String()
	fmt.Println("ipStr:", ipStr)
	defer func() {
		fmt.Println("disconnected:", ipStr)
		conn.Close()
		deleteConn(conn)
		boradcastMessage(ipStr + "离开了房间" + "\n")
	}()
	reader := bufio.NewReader(conn)
	for {
		// 读取直到输入中第一次发生 '\n'
		message, err := reader.ReadString('\n')
		//因为按强制退出的时候，他就先发送换行，然后在结束
		if message == "\n" {
			return
		}
		message = ipStr+"说："+message
		if err != nil {
			fmt.Println("topPipe:",err)
			return
		}

		// 广播消息
		fmt.Println(ipStr,"说：",message)
		err = boradcastMessage(message)
		if err!=nil{
			fmt.Println(err)
			return
		}
	}
}

// 广播数据
func boradcastMessage(message string) error {
	b := []byte(message)
	for i := 0; i < len(connSlice); i++ {
		fmt.Println(connSlice[i])
		_, err := connSlice[i].Write(b)
		if err != nil {
			fmt.Println("发送给", connSlice[i].RemoteAddr().String(), "数据失败"+err.Error())
			continue
		}
	}
	return nil
}

// 移除已经关闭的客户端
func deleteConn(conn *net.TCPConn) error {
	if conn == nil {
		fmt.Println("conn is nil")
		return errors.New("conn is nil")
	}
	for i := 0; i < len(connSlice); i++ {
		if connSlice[i] == conn {
			connSlice = append(connSlice[:i], connSlice[i+1:]...)
			break
		}
	}
	return nil
}
