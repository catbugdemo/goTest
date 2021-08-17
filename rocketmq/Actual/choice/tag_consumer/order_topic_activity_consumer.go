package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"main.go/rocketmq/constant"
	"time"
)

type OrderInfo struct {

}

// 只要用户下单并成功支付，就发放一张优惠券
func main() {
	// 该admin可以进行优化,可以不用写， broker-a.properties
	// defaultTopicQueueNums=4 #在发送消息时，自动创建服务器不存在的topic，默认创建的队列数
	// https://www.huaweicloud.com/articles/47e7a12d7c922d12666a970f5e6fc6be.html
	testAdmin, err := admin.NewAdmin(admin.WithResolver(primitive.NewPassthroughResolver(constant.Namesrv)))
	if err != nil {
		panic(err)
	}
	// 判断是否存在，存在则创建节点地址
	if err = testAdmin.CreateTopic(
		context.Background(),
		admin.WithTopicCreate(constant.Topic),
		admin.WithBrokerAddrCreate("118.89.121.211:10911"),
	);err!=nil{
		panic(err)
	}

	c, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName(constant.TestGroup),
		consumer.WithNameServer(constant.Namesrv),
		// 指定发送失败时的重试时间
		consumer.WithRetry(2),
	)
	if err != nil {
		panic(err)
	}
	selector := consumer.MessageSelector{
		Type:       consumer.TAG,
		Expression: "TagA || TagC",
	}

	// 先写消费信息，再开始消费
	if err = c.Subscribe(constant.Topic, selector, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, msg := range msgs {
			fmt.Println(msg.Body)
		}
		return consumer.ConsumeSuccess, nil
	});err!=nil{
		panic(err)
	}

	if err = c.Start();err!=nil{
		panic(err)
	}
	time.Sleep(time.Hour)
	if err = c.Shutdown();err!=nil{
		panic(err)
	}
}

// MsgId 消息发送者在消息发送时会首先在客户端生成的全局唯一