package test

import (
	"fmt"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"testing"
	"time"
)

func TestFibonacci(t *testing.T) {
	Convey("test Fibonacci", t, func() {
		So(Fibonacci(2),ShouldEqual,1)
	})
}

func Fibonacci(num int) int {
	n1, n2 := 0, 1
	if num == 1 {
		return n1
	}
	if num == 2 {
		return n2
	}
	for i := 0; i < num; i++ {
		sum := n1 + n2
		n1 = n2
		n2 = sum
	}
	return n2
}

func TestUnix(t *testing.T) {
	t.Run("unix", func(t *testing.T) {
		i := time.Now().Unix() % 30
		fmt.Println(i)
	})
}

func TestA(t *testing.T) {
	t.Run("err", func(t *testing.T) {
		log.Println(errors.Wrap(errors.New("错误"),"一个"))
	})
}
