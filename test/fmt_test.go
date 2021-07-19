package test

import (
	"fmt"
	"testing"
)

func TestFmt(t *testing.T) {
	t.Run("fmt", func(t *testing.T) {
		type Test struct {
			Name string
		}
		ints := Test{
			Name: "Tom",
		}
		// 默认格式输出
		fmt.Printf("%v", ints) //  {Tom}
		fmt.Println()

		// 输出结构体时会添加字段名
		fmt.Printf("%+v", ints) // {Name:Tom}
		fmt.Println()

		// %#v 的输出形式
		fmt.Printf("%#v", ints) // test.Test{Name:"Tom"}
		fmt.Println()

		// %T	值的类型的Go语法表示
		fmt.Printf("%T", ints) // test.Test
		fmt.Println()

		// %% 百分号
		fmt.Printf("%%") // %

	})
}

func TestMake(t *testing.T) {
	t.Run("not make", func(t *testing.T) {
		ints := make([][]int, 2, 10)
		i := append(ints, []int{1, 2})
		fmt.Println(i)
	})
}


// 格式化输出json ()
func formatJSON()  {

}