package poolTest

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExample(t *testing.T) {
	t.Run("pool", func(t *testing.T) {
		pool, err := NewPool(20)
		assert.Nil(t, err)
		defer pool.Close()

		for i := 0; i < 20; i++ {
			pool.Put(&Task{
				Handler: func(v ...interface{}) {
					fmt.Println(v)
				},
				Params: []interface{}{i},
			})
		}
	})
}
