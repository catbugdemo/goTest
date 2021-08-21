package test

import (
	"encoding/json"
	"sync"
	"testing"
)

type Student struct {
	Name string
	Age  int
}

var buf, _ = json.Marshal(Student{Name: "Geektutu", Age: 25})


func BenchmarkMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		st := &Student{}
		json.Unmarshal(buf,st)
	}
}

var commonPool = sync.Pool{
	New: func() interface{} {
		return new(Student)
	},
}

func BenchmarkMarshalWithPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		st := commonPool.Get().(*Student)
		json.Unmarshal(buf,st)
		commonPool.Put(st)
	}
}
