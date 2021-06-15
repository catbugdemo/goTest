package poolTest

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPool(t *testing.T) {
	t.Run("pool", func(t *testing.T) {
		pool, err := NewPool(20)
		assert.Nil(t, err)

		for i := 0; i < 20; i++ {
			pool.Put(&Task{
				Handler: func(v ...interface{}) {
					fmt.Println(v)
				},
				Params: []interface{}{i},
			})
		}
		pool.Close()
		fmt.Println(pool.status)
	})
}