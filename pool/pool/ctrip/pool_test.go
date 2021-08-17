package ctrip

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	t.Run("pool", func(t *testing.T) {
		//创建任务池
		pool, e := NewPool(10)
		if e !=nil{
			panic(e)
		}

		for i := 0; i < 1; i++ {
			pool.Put(&Task{
				Handler: func(v ...interface{}) {
					fmt.Println(v)
				},
				Params: []interface{}{i},
			})
		}

		time.Sleep(1e9)
		pool.Close()
		fmt.Println(pool.status)
	})
}

var wg = sync.WaitGroup{}

var sum int64

func demoTask(v ...interface{}) {
	defer wg.Done()
	for i := 0; i < 100; i++ {
		atomic.AddInt64(&sum, 1)
	}
}

func Test100W(t *testing.T) {
	var runTimes = 1000000

	t.Run("test", func(t *testing.T) {
		for i := 0; i < runTimes; i++ {
			wg.Add(1)
			go demoTask()
		}
	})

	t.Run("pool", func(t *testing.T) {
		pool, err := NewPool(20)
		assert.Nil(t, err)

		task := &Task{
			Handler: demoTask,
		}

		for i := 0; i < runTimes; i++ {
			wg.Add(1)
			pool.Put(task)
		}
	})
}
