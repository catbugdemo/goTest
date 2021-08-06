package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"os"
	"strconv"
)

func main() {
	p, _ := rocketmq.NewProducer(
		// 设置  nameSrvAddr
		// nameSrvAddr 是 Topic 路由注册中心
		producer.WithNameServer([]string{"118.89.121.211:9876"}),
		// 设置 Retry 重连次数
		producer.WithRetry(2),
		// 设置 Group
		producer.WithGroupName("testGroup"),
	)
	// 开始连接
	err := p.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}

	// 设置节点名称
	topic := "Topic-test"
	// 循坏发送信息 (同步发送)
	for i := 0; i < 10; i++ {
		msg := &primitive.Message{
			Topic: topic,
			Body:  []byte("Hello RocketMQ Go Client" + strconv.Itoa(i)),
		}
		// 发送信息
		res, err := p.SendSync(context.Background(),msg)
		if err != nil {
			fmt.Printf("send message error:%s\n",err)
		}else {
			fmt.Printf("send message success: result=%s\n",res.String())
		}
	}
/*
	// 异步发送
	size := 10
	var wg sync.WaitGroup
	for i := 0; i < size; i++ {
		wg.Wait()
		msg := &primitive.Message{
			Topic: "test",
			Body:  []byte(fmt.Sprintf("%d: 你好，Go 客户端，这是异步消息", i)),
		}
		p.SendAsync(context.Background(), func(ctx context.Context, result *primitive.SendResult, err error) {
			defer func() {
				wg.Done()
			}()
			if err != nil {
				fmt.Printf("消息发送成功：%s\n",err)
			}
		},msg)
	}
	wg.Wait()
*/

	// 关闭生产者
	err = p.Shutdown()
	if err != nil {
		fmt.Printf("shutdown producer error:%s",err.Error())
	}
}


