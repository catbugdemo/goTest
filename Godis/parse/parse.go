package parse

import (
	"io"
	"main.go/Godis/interface/redis"
)

//Payload 包括一个func []byte Data 和 error
type Payload struct {
	Data redis.Reply
	Err  error
}

//ParseStream 通过io.Reader 读取数据并将结果通过 channel将结果返回给调用者
//流式处理的接口适合供客户端/服务端使用
func ParseStream(reader io.Reader) <-chan *Payload {
	//创建结构体
	ch := make(chan *Payload)
	//运行
	go parse0(reader, ch)
	return ch
}

// parse0 核心流程
func parse0(reader io.Reader, ch chan<- *Payload) {
	//初始化读取状态
	readingMultLine := false
	expectedArgsCount := 0
	var args [][]byte
	var bulkLen int64 // 数组长度 又称 Bulk String数组
	for {
		//RESP是以行为单位的
		//行分为简单字符串和二进制安全的BulkString,我们需要封装一个 readLine函数来兼容
		line, e := readLine(reader, bulkLen)
		if e != nil {
			//处理错误
		}

		//对刚刚读取的进行解析
		//间断将Reply分为两类
		//单行 StatusReply(状态返回),IntReply(整数返回),ErrorReply(错误)


	}

}

// readLine 读取行数
func readLine(reader io.Reader, bulkLen int64) ([]byte,error) {

}

