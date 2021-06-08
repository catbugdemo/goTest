package main

import (
	"main.go/redis/pkg/setting"
)

func main() {
	//初始化全局配置
	setting.InitSetting()

	//初始化缓存
	setting.InitRedisPool()
}
