package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"strconv"
	"time"
)

// 消息有序，顺序消费(FIFO - 先进先出)

//设置结构体
type OrderStep struct {
	OrderId uint   `json:"order_id"`
	Desc    string `json:"desc"`
}

type QueueSelect struct {
}

// 获取消息队列中的 msg , mqs
func (o *QueueSelect) Select(msg *primitive.Message, mqs []*primitive.MessageQueue) *primitive.MessageQueue {
	// 根据订单id 发送 queue
	orderID, _ := strconv.Atoi(msg.GetProperty("orderID"))
	// 计算队列索引
	index := orderID % len(mqs)
	return mqs[index]
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

func main() {
	groupName := "testGroup"
	nameServe := []string{"118.89.121.211:9876"}
	// 实例化队列选择器
	qs := &QueueSelect{}



	// 创建生产者
	p, _ := rocketmq.NewProducer(
		producer.WithGroupName(groupName),
		producer.WithRetry(2),
		producer.WithNameServer(nameServe),
		producer.WithQueueSelector(qs),
	)
	// 开始与结束
	p.Start()

	// 创建标签
	tags := []string{"TagA","TagB","TagD"}

	// 订单列表
	orders := buildOrders()
	for i := 0; i < 10; i++ {
		msg := &primitive.Message{
			Topic: "Topic-Test-Demo",
			Body:  []byte(time.Now().Format("2006-01-02 15:04:05") + "Hello rocketmq" + fmt.Sprintf("OrderStep{orderID=%d, desc=%s}", orders[i].OrderId, orders[i].Desc)),
		}
		// 将订单 id 放到消息属性中
		msg.WithProperty("orderID",fmt.Sprint(orders[i].OrderId))
		// 消息添加标签 (随机标签)
		msg.WithTag(tags[i%len(tags)])
		// 消息添加 key
		msg.WithKeys([]string{"Key-"+ fmt.Sprint(orders[i].OrderId)})
		// 延迟3级，即延迟 10s
		msg.WithDelayTimeLevel(3)

		// 消息同步发送
		sync, err := p.SendSync(context.Background(), msg)
		if err != nil {
			panic(err)
		}
		fmt.Printf("消息发送成功,结果状态：%d,队列id:%d\n",sync.Status,sync.MessageQueue.QueueId)
	}

	if err := p.Shutdown(); err != nil {
		panic(err)
	}
}

// 参考：https://cloud.tencent.com/developer/article/1581368
// https://blog.thepoy.cn/2021/05/13/RocketMQ-4.8.0-%E5%85%A5%E9%97%A8%E8%AE%B0%E5%BD%95-%E4%BA%8C-Golang-%E5%AE%A2%E6%88%B7%E7%AB%AF.html#2-%E9%A1%BA%E5%BA%8F%E6%B6%88%E6%81%AF
// 延迟消费参考：https://zhuanlan.zhihu.com/p/379547780

// 1.发送新消息可以延时吗
// 2.能异步发送消息吗
// 3.发送消息能够拦截吗
// 4.能对 Tag 进行简单的逻辑判断吗
// 5.能先启动 consumer 再启动 producer吗
//   如果消费者先启动，可能会存在未创建 Topic 节点 而进行连接
// 	(1) 可以通过 admin 先创建节点 topic,再进行消费
//	(2) 可以通过 broker 中的 .properties 配置节点，自动创建，但是生产环境不推荐
// 6.消息发送时如果使用了MessageQueueSelector,那消息发送的重试机制将失效
//		即 RocketMQ 客户端并不会重试，消息发送的高可用需要由业务方来保证，
//		典型的就是消息发送失败后存在数据库中，然后定时调度，最终将消息发送到MQ中
// 7.rocketmq 批量消费容量是否有一个限度，限度是多少？ 4mb
// 8.消费者或生产者的第一次连接是连接主服务器还是从服务器
// 9.rocketmq的提交与回滚
// 10.消费者消息过滤的3种方式：Tag ,SQL语句，Property

// 消息发送时常见错误与解决方案
// the topic = route info not found, it may not exist
// 创建节点

// System busy : too many requests and system thread pool busy
// Broker busy:
//		PageCache 繁忙 锁超过1s ，开启 transientStroePoolEnable=true 机制 (开启内存锁)
//		异步批量提交然后 PageCache 进行 同步，异步刷盘