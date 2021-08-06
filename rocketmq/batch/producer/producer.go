package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"strconv"
)

// 批量发送消息
func main() {
	// 设置新节点
	topic := "BatchTest"
	nameServer := []string{"118.89.121.211:9876"}

	/*
		创建一个 topic
	*/
	testAdmin, err := admin.NewAdmin(admin.WithResolver(primitive.NewPassthroughResolver(nameServer)))
	if err != nil {
		panic(err)
	}

	// 创建节点
	if err = testAdmin.CreateTopic(
		context.Background(),
		// 创建 topic
		admin.WithTopicCreate(topic),
		// 创建节点的地址
		admin.WithBrokerAddrCreate("118.89.121.211:10911"),
	); err != nil {
		panic(err)
	}

	// 创建生产者
	p, _ := rocketmq.NewProducer(
		producer.WithCreateTopicKey(topic),
		producer.WithNameServer(nameServer),
		producer.WithRetry(2),
	)

	if err = p.Start(); err != nil {
		panic(err)
	}
	defer func() {
		if err = p.Shutdown(); err != nil {
			panic(err)
		}
	}()

	var msgs []*primitive.Message
	for i := 0; i < 10; i++ {
		msg := &primitive.Message{
			Topic: topic,
			Body:  []byte([]byte("你好， RocketMQ Go 客户端！ num: " + strconv.Itoa(i))),
		}
		// 设置 Tag
		msg.WithTag("TagA")
		// 相当于数据库中的 id
		msg.WithKeys([]string{"OrderID00" + strconv.Itoa(i+1)})
		// 添加筛选属性
		msg.WithProperty("a", strconv.Itoa(i))
		msgs = append(msgs, msg)
	}

	sync, err := p.SendSync(context.Background(), msgs...)
	if err != nil {
		panic(err)
	}
	fmt.Printf("消息发送成功: result=%s\n",sync.String())
}

// 添加简单逻辑
// echo "enablePropertyFilter = true" >> conf/broker.conf