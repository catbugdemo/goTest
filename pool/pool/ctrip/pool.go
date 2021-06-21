package ctrip

import (
	"errors"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

//每个结构体都需要New 来创建他自身的结构体
//

//任务定义 需要执行函数，以及函数要传的参数
type Task struct {
	Handler func(v ...interface{})
	Params  []interface{}
}

//定义连接池
type Pool struct {
	//池的容量
	capacity uint64
	//数量
	runningWorkers uint64
	//任务池状态
	status int64
	//任务队列
	chTask chan *Task
	//提供加锁接口
	sync.Mutex
	//提供可定制的panic handler
	PanicHandler func(interface{})
}

//任务池构造函数
var ErrInvalidPoolCap = errors.New("invalid pool cap")
var ErrPoolAlreadyClosed = errors.New("pool already closed")

const (
	STOPED = iota
	RUNNING
)

// NewPool 任务池构造函数
func NewPool(capacity uint64) (*Pool, error) {
	if capacity <= 0 {
		return nil, ErrInvalidPoolCap
	}
	// 配置连接池: 容量:capacity , 状态：status 运行状态
	return &Pool{
		capacity: capacity,
		status:   RUNNING,
		//初始化任务队列，配置队列长度为容量
		chTask: make(chan *Task, capacity),
	}, nil
}

// run
// 作为启动worker的方法
func (p *Pool) run() {
	// 运行中任务+1
	p.incRunning()

	go func() {
		defer func() {
			//worker结束，运行中的任务-1
			p.decRunning()
			//恢复panic
			if r := recover();r!=nil{
				if p.PanicHandler !=nil{
					p.PanicHandler(r)
				}else {
					// 默认处理
					log.Printf("Worker panic: %s\n",r)
				}
			}
			// worker 退出时检测是否有可运行的worker
			p.checkWorker()
		}()

		for {
			//阻塞等待任务，结束信号到来
			select {
			//从channel中消费任务
			case task,ok:= <-p.chTask:
				// 如果channel被关闭，结束worker运行
				if !ok {
					return
				}
				//执行任务
				task.Handler(task.Params...)
			}
		}
	}()
}

// GetCap() 获取容量
func (p *Pool) GetCap() uint64 {
	return p.capacity
}

// 对runningWorkers的操作进行封装
func (p *Pool) incRunning() {
	atomic.AddUint64(&p.runningWorkers,1)
}

func (p *Pool) decRunning()  {
	atomic.AddUint64(&p.runningWorkers,^uint64(0))
}

func (p *Pool) GetRunningWorkers() uint64 {
	return atomic.LoadUint64(&p.runningWorkers)
}

// setStatus 设置status状态
func (p *Pool) setStatus(status int64) bool {
	p.Lock()
	defer p.Lock()

	if p.status == status{
		return false
	}

	p.status = status
	return true
}

func (p *Pool) Put(task *Task) error {
	// 加锁防止启动多个worker
	p.Lock()
	defer p.Unlock()

	//如果任务池处于关闭状态，再put任务会返回ErrPool
	if p.status == STOPED {
		return ErrPoolAlreadyClosed
	}

	// 任务池未满，创建worker
	if p.GetRunningWorkers() <p.GetCap(){
		p.run()
	}

	// send task,将任务通过通道
	if p.status == RUNNING {
		p.chTask <- task
	}
	return nil
}

func (p *Pool) Close() {
	//设置 status为已停止状态
	p.setStatus(STOPED)

	//阻塞等待所有任务被worker,消费
	for len(p.chTask) > 0 {
		time.Sleep(1e6)
	}

	//关闭任务队列
	close(p.chTask)
}

func (p *Pool) checkWorker() {
	p.Lock()
	defer p.Lock()

	//当没有worker 且有任务存在时
	// 运行一个 worker 消费任务
	if p.runningWorkers == 0 && len(p.chTask) > 0{
		p.run()
	}
}