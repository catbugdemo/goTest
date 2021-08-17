package test

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

// LinkRedisWayOne 测试redis连接第一种方法
func LinkRedisWayOne() {
	dial, err := redis.Dial("tcp", "118.89.121.211:6379")
	if err != nil {
		panic(err)
	}
	fmt.Println("redis link one way success")
	defer dial.Close()

}

// LinkRedisWayTwo 测试redis连接第二种方法
func LinkRedisWayTwo() {
	url, err := redis.DialURL("redis://118.89.121.211:6379")
	if err != nil {
		panic(err)
	}
	fmt.Println("redis link two way success")
	defer url.Close()
}