package poolTest

import (
	"errors"
	"sync"
)

const (
	STOPED = iota
	RUNNING
)

var ErrInvalidPoolCap = errors.New("invalid pool cap")
var ErrPoolAlreadyClosed = errors.New("pool already closed")

// Task 任务
type Task struct {
	//传入函数
	Handler func(v ...interface{})
	//需要传入的v
	Params  []interface{}
}

// Pool 池
type Pool struct {
	//最大容量
	MaxCap int
	//正在进行的任务数量
	active int
	//池的状态 RUNNING OR STOPED
	status int
	//任务通道
	chTask chan *Task
	//锁
	mu     sync.Mutex
	//用于等待一组线程的结束
	wg     sync.WaitGroup
}

//创建
func NewPool(maxCap int) (*Pool, error) {
	if maxCap <= 0 {
		return nil, ErrInvalidPoolCap
	}

	return &Pool{
		MaxCap: maxCap,
		status: RUNNING,
		chTask: make(chan *Task, maxCap),
	}, nil
}

// Put 就绪
// 1.加锁
// 2.判断池是否为运行状态
func (p *Pool) Put(task *Task) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.status == STOPED {
		return ErrPoolAlreadyClosed
	}

	//放入消息队列中
	p.chTask <- task

	p.wg.Add(1)
	//运行线程
	if p.active < p.MaxCap {
		p.run()
	}
	return nil
}

// run 运行
func (p *Pool) run() {
	p.active++
	go p.worker()
}

//worker 阻塞
func (p *Pool) worker() {
	defer func() {
		p.active--
		p.wg.Done()
	}()

	for {
		select {
		case task, ok := <-p.chTask:
			if !ok {
				return
			}
			task.Handler(task.Params...)
		}
	}
}

// Close 结束
// 1.设置状态为结束
// 2.等待所有线程执行完毕
// 3.关闭线程
func (p *Pool) Close() {
	p.status = STOPED
	close(p.chTask)
	p.wg.Wait()
}
