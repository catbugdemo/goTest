package pool

import (
	"errors"
	"sync"
)

// 任务 Struct
type Task struct {
	// 需要执行的函数
	Handle func(v ...interface{})
	Params []interface{}
}

const (
	STOPED = iota
	RUNNING
)

type Pool struct {
	// 最大容量
	MaxCap int
	// 正在进行的任务数量
	active int
	// 池的状态
	status int
	// 任务通道
	chTask chan *Task
	// 锁
	mu sync.Mutex
	// 等待线程
	wg sync.WaitGroup
}

// 新建 就绪 运行 阻塞 结束

// 创建线程最大数量
func NewPool(maxCap int) (*Pool, error) {
	if maxCap < 0 {
		return nil, errors.New("maxCap 小于0")
	}

	return &Pool{
		MaxCap: maxCap,
		status: RUNNING,
		chTask: make(chan *Task,maxCap),
	}, nil
}

// Put 就绪 将任务放入队列中
func (p *Pool) Put(task *Task) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.status == STOPED {
		return errors.New("没启动")
	}

	p.chTask <- task

	p.wg.Add(1)
	// 运行线程,活动数小于 最大容量
	if p.active < p.MaxCap {
		p.run()
	}

	return nil
}

// 运行
func (p *Pool) run() {
	p.active++
	go p.worker()
}

// 阻塞
func (p *Pool) worker() {
	defer func() {
		p.active--
		p.wg.Done()
	}()

	// 起到阻塞作用
	for {
		select {
		case task, ok := <-p.chTask:
			if !ok {
				return
			}
			task.Handle(task.Params...)
		}
	}
}

// 关闭
func (p *Pool) Close()  {
	p.status = STOPED
	close(p.chTask)
	p.wg.Wait()
}