package poolMyself

import (
	"errors"
	"sync"
)

// 使用代理模式:
//	Task---Subject
//	Pool---Proxy
//
// 线程状态:
// 创建(new),就绪(runnable),运行(run),阻塞(block),结束(dead)
// 死锁:
// 破坏互斥条件
// 破坏占有和等待条件
// 破坏不剥夺条件
// 破坏循环等待条件
// https://zhuanlan.zhihu.com/p/61221667
// https://blog.csdn.net/hd12370/article/details/82814348

// 1.创建一个 协程池需要的Task属性
// 包括:
//	func 函数
//	自身状态
//
// 2.创建一个协程池需要的属性
// 包括:
//	MaxCap --- 最大容量
//	chTask --- 任务
//  Status --- 自身状态
//  active --- 活动数量
//	sync.Mutex --- 锁

// 3.都需要的New一个池，一个池最重要的就是
//		开始 --- run
//		关闭 --- close
// 4.任务自身需要运行的内容 run
// 5.池进行的一个代理

var ErrInvalidPoolCap = errors.New("invalid pool cap")
var ErrPoolAlreadyClosed = errors.New("pool already closed")

// 运行状态
const (
	STOPED = iota
	RUNNING
)

// Task 具体对象
type Task struct {
	Handler func(v int)
	Params  int
}

// Pool 池
type Pool struct {
	MaxCap int
	status int
	chTask chan *Task
	active int
	mu     sync.Mutex
	wg     sync.WaitGroup
}

// NewPool 新建一个池 --创建
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

// worker 等待
// 1.阻塞任务
// 2.运行任务
// 3.等待任务
func (p *Pool) worker() {
	for {
		defer func() {
			p.active--
		}()
		select {
		//释放线程资源
		case task, ok := <-p.chTask:
			if !ok {
				return
			}
			//运行任务
			task.Handler(task.Params)
		}
	}
}

// run 运行
// 1.活动状态+1
// 2.运行 -- 等待
func (p *Pool) run() {
	// 1.活动状态+1
	p.active++
	go p.worker()
}

// Put 就绪 --- 将需要的任务放入池中
// 1.判断状态
// 2.判断线程池容量 --- 大于阻塞
// 3.添加任务
// 4.转为运行态
func (p *Pool) Put(task *Task) error {
	// 1.加锁 --- 破坏循环等待条件
	p.mu.Lock()
	defer p.mu.Unlock()

	// 1.判断状态
	if p.status == STOPED {
		return ErrPoolAlreadyClosed
	}

	// 2.判断线程池容量
	if p.active < p.MaxCap {
		p.run()
	}

	// 3.添加任务
	p.chTask <- task

	return nil
}

//Close 关闭
// 1.设置池状态为结束状态
// 2.判断是否还有线程
// 3.关闭线程
func (p *Pool) Close() {
	p.status = STOPED

	for len(p.chTask) > 0 {
		p.run()
	}

	close(p.chTask)
}
