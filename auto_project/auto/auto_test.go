package auto

import (
	"fmt"
	"reflect"
	"testing"
)

type T struct {
	//omitempty 省略空
	ID int `json:"id"`
	Name string`json:"name"`
}

func TestInterface(t *testing.T) {
	t.Run("interface", func(t *testing.T) {
		vType := reflect.TypeOf(T{})
		fmt.Println(vType.String())
	})
}
