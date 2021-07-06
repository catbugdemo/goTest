package main

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		str := "AaDaB"
		for i := 0; i < len(str); i++ {
			if str[i] >= 65 && str[i] <= 90 {
				str = str[0:i] + "_" + str[i:]
				i++
			}
		}
		fmt.Println(str)
	})
}

func TestA(t *testing.T) {
	type A struct {
		M interface{}
		S []string
		I time.Time
	}
	a := A{
		M: "A",
		S: []string{"1","2"},
	}
	of := reflect.TypeOf(a)
	s := of.Field(2).Type.String()
	sprintf := fmt.Sprintf("%s", s)
	fmt.Println(sprintf)
}