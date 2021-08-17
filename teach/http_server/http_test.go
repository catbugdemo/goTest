package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

// 测试端口 5000次 每次循环 244450 ns
func BenchmarkName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		get, _ := http.Get("http://127.0.0.1:8082/self")
		all, _ := ioutil.ReadAll(get.Body)
		fmt.Println(string(all))
	}
}

func TestHttp(t *testing.T) {
	get, e := http.Get("http://127.0.0.1:8082/self")
	assert.Nil(t, e)

	all, e := ioutil.ReadAll(get.Body)
	assert.Nil(t, e)
	fmt.Println(string(all))
}
