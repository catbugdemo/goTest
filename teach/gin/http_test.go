package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

// 测试 gin
// 测试 5000 循环 ， 每次循环请求耗时 221624 纳秒
func BenchmarkGin(b *testing.B) {
	for i := 0; i < b.N; i++ {

		get, e := http.Get("http://127.0.0.1:8081/test")
		if e != nil {
			panic(e)
		}
		all, e := ioutil.ReadAll(get.Body)
		if e != nil {
			panic(e)
		}
		fmt.Println(string(all))
	}
}

func TestHttp(t *testing.T) {
	get, e := http.Get("http://127.0.0.1:8081/gin")
	assert.Nil(t, e)
	all, e := ioutil.ReadAll(get.Body)
	assert.Nil(t, e)
	fmt.Println(string(all))
}
