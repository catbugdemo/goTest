package errorz

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrintFormt(t *testing.T) {
	t.Run("集成测试：Pack", func(t *testing.T) {
		err := Pack(errors.New("test"))

		fmt.Print(err)
		assert.NotNil(t, err)
	})

	t.Run("单元测试:PrintFormat", func(t *testing.T) {
		err := errors.New("Error test")
		format := PrintFormat("test", 12, err.Error())

		fmt.Println(format)
		assert.NotNil(t, format)
	})
}
