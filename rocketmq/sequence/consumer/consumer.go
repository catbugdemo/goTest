package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"time"
)

func main() {
	groupName := "testGroup"
	nameServer := []string{"118.89.121.211:9876"}

	// 创建消费者
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName(groupName),
		consumer.WithRetry(2),
		consumer.WithNameServer(nameServer),
		// 设置 consumer 第一次启动是从队列头部开始还是队列尾部开始消费
		// 如果非第一次启动，那么按照上次消费的位置继续消费
		consumer.WithConsumeFromWhere(consumer.ConsumeFromFirstOffset),
		// 设置为顺序消费
		consumer.WithConsumerOrder(true),
		// 延时消费
	)

	err := c.Subscribe(
		"123",
		consumer.MessageSelector{Expression: "TagA || TagB || TagD"},
		// consumer.MessageSeletor{Expression: "a between 110 and 130"}
		func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
			// 获取内容信息
			orderlyCtx, _ := primitive.GetOrderlyCtx(ctx)
			fmt.Println("顺序消费的上下文:", orderlyCtx)
			fmt.Println("顺序消费的回调:", msgs)
			return consumer.ConsumeSuccess, nil
		},
	)
	if err != nil {
		panic(err)
	}
	// 启动必须在订阅之后
	if err = c.Start();err!=nil{
		panic(err)
	}
	fmt.Println("开始消费")

	// 模拟业务处理时间
	time.Sleep(time.Hour)
	if err = c.Shutdown();err!=nil{
		panic(err)
	}
}

// 延迟消费时间一般可以通过 producer 设置响应时间
// 延迟消息使用场景:
//		比如电商，提交了一个订单就可以发送一个延时新消息，1小时候去检查这个订单的状态
//		如果还未付款，就取消订单释放库存

// namespace 是先于 broker 启动还是后于
// broker 部署的几种状态 单机部署，多master，双主双从，Dledger (双主双从+ raft [redis 哨兵模式])
// consumer 先行启动时，如果设置一个不存在的 Topic 会直接 panic 报错，可以用admin，
// 也可以通过 broker-a.properties 中的 defaultTopicQueueNums= 自行创建
// namesrv 之间相互通信吗 ， 不通信

// 消息消费积压问题
//	消息积压的问题大概率是消费端 consumer 的问题
//	消息发送遇到的问题更有可能是 broker端的问题

// rocketmq不保证消费的幂等性，但保证至少消费一次