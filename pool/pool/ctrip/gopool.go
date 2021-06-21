package ctrip

import (
	"fmt"
	"time"
)

type GoTask struct {
	// 一个无参构造方法
	Handler func() error
}

//NewTask 创建一个新的任务
func NewTask(handler func() error) *GoTask {
	goTask := GoTask{
		Handler: handler,
	}
	return &goTask
}

//执行Task任务的方法
func (t *GoTask) Execute() {
	t.Handler()
}

//定义池的结构体
type GoPool struct {
	//入口
	EntryChannel chan *GoTask
	//最大连接数
	capacity int
	//就绪队列
	JobsChannel chan *GoTask
}

// NewGoPool 新建一个连接池
func NewGoPool(capacity int) (*GoPool, error) {
	if capacity <= 0 {
		return nil, ErrInvalidPoolCap
	}
	return &GoPool{
		EntryChannel: make(chan *GoTask, capacity),
		capacity:     capacity,
		JobsChannel:  make(chan *GoTask, capacity),
	}, nil
}

// 协成池创建一个worker并开始工作
func (p *GoPool) worker(id int) {
	for task := range p.JobsChannel {
		task.Execute()
		fmt.Println("worker ID", id, "执行完毕任务")
	}
}

func (p *GoPool) Run() {
	//固定数量
	for i := 0; i < p.capacity; i++ {
		go p.worker(i)
	}

	//从 EntryChannel 协程池入口取外界传递过来的任务
	for task := range p.EntryChannel {
		p.JobsChannel <- task
	}

	//执行完后关闭JobsChannel，该方法使用官方的内建函数close关闭信道
	close(p.JobsChannel)
	//关闭EntryChannel
	close(p.EntryChannel)
}

//主函数
func main() {
	task := NewTask(func() error {
		fmt.Println(time.Now())
		return nil
	})

	pool, err := NewGoPool(3)
	if err != nil{
		panic("fail to new pool")
	}

	go func() {
		for {
			pool.EntryChannel <- task
		}
	}()

	//启动协程池
	pool.Run()
}
