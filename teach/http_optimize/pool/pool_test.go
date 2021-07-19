package pool

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPool(t *testing.T) {
	t.Run("pool", func(t *testing.T) {
		pool, e := NewPool(20)
		assert.Nil(t, e)
		defer pool.Close()

		for i := 0; i < 20; i++ {
			pool.Put(&Task{
				Handle: func(v ...interface{}) {
					fmt.Println(v)
				},
				Params: []interface{}{i},
			})
		}
	})
}
