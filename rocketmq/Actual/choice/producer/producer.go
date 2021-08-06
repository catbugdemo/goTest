package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"log"
	"main.go/rocketmq/constant"
	"strconv"
	"time"
)

// 队列选择机制
// 订单系统通可以允许用户更新新订单的信息，
// 并且订单有其流转的生命周期，
// 如待付款，已支付，卖家已发货，买家已收货等等

// 消息发送时如果使用了 MessageQueueSelector,那消息发送的重试机制将失效
// 即 RocketMQ 客户端并不会重试，消息发送的高可用需要由业务方来保证
// 典型的就是消息发送失败后存在数据库中，然后定时调度，最终将新消息发送到MQ

// RocketMQ 支持队列级别的顺序消费(FIFO) 我们只需要在学校学习发送的时候
// 如果将同一个订单号的不同消息发送到同一个队列，
// 这样在消费的时候，我们就能按照顺序进行处理

type querySelect struct {
}

// Select 用来选择 MessageQueue
// 我们相同订单的创建，修改，删除，都有一个共同的 OrderId
// 将相同的 OrderId 的创建修改删除都按照顺序队列存储在同一个 Topic 中
// msg - 当前处理的原始信息, mqs 节点信息
func (q *querySelect) Select(msg *primitive.Message, mqs []*primitive.MessageQueue) *primitive.MessageQueue {
	// 根据相同的 OrderId 选择 Topic
	orderId, _ := strconv.Atoi(msg.GetProperty("OrderId"))
	if orderId == 0 {
		return nil
	}
	// 通过取余插入到相同的 Broker中
	index := orderId % orderId
	return mqs[index]
}

//设置结构体
type OrderStep struct {
	OrderId uint   `json:"order_id"`
	Desc    string `json:"desc"`
}

// 生成模拟订单
func buildOrders() []*OrderStep {
	orderSteps := make([]*OrderStep, 0, 20)
	var strl = []string{"创建", "付款", "推送", "完成"}
	for i := 0; i < 4; i++ {
		for j := 0; j < 3; j++ {
			o := new(OrderStep)
			o.OrderId = 15103111030 + uint(j)
			o.Desc = strl[i]
			orderSteps = append(orderSteps, o)
		}
	}
	return orderSteps
}

// 单独对于 producer
func main() {
	qs := &querySelect{}

	p, err := rocketmq.NewProducer(
		// 需要的默认值 group,namsrv,witRetry,
		producer.WithGroupName(constant.TestGroup),
		producer.WithNameServer(constant.Namesrv),
		producer.WithRetry(2),
		// 自定义默认选择器
		producer.WithQueueSelector(qs),
	)
	if err != nil {
		panic(err)
	}

	// 开始与结束
	if err = p.Start();err!=nil{
		panic(err)
	}
	defer func() {
		if err = p.Shutdown();err!=nil{
			panic(err)
		}
	}()
	tags := []string{"TagA","TagB","TagD"}

	orders := buildOrders()
	// 发送同步消息
	for i := 0; i < 10; i++ {
		msg := &primitive.Message{
			Topic: constant.Topic,
			Body:  []byte(time.Now().Format("2006-01-02 15:04:05") + "Hello rocketmq" + fmt.Sprintf("OrderStep{orderID=%d, desc=%s}", orders[i].OrderId, orders[i].Desc)),
		}
		msg.WithProperty("OrderId",fmt.Sprint(orders[i].OrderId))
		// Tag 实现对不同子主题的不同消费逻辑，实现更好的扩展性
		msg.WithTag(tags[i%len(tags)])
		// 业务标识 Key
		msg.WithKeys([]string{"Key-"+ fmt.Sprint(orders[i].OrderId)})

		sync, err := p.SendSync(context.Background(), msg)
		if err != nil {
			log.Println("Fail to send:",orders[i].OrderId," ",msg)
			// 因为自定义了消息队列选择，所以不会重试，可以存入数据库中进行重试
			// 数据库操作
		}
		fmt.Printf("消息发送成功,结果状态：%d,队列id:%d，队列信息:%d %s\n",sync.Status,sync.MessageQueue.QueueId,orders[i].OrderId,orders[i].Desc)
	}
}
