package main

import (
	"fmt"
	"log"
	"main.go/auto_project/data_scrpit/code"
	"os"
	"strconv"
)

func main() {
	// 数据存储函数
	strs := make([]string, 0, 4)

	for i, arg := range os.Args {
		if arg == "-help" {
			fmt.Println("如果想生成 hdfs to pg 请输入2个参数如: ./main table_name isTest (table_name 是表名， isTest 是 是否是测试服 true or false）")
			fmt.Println("如果想生成 pg to hdfs 请输入3个参数如:./main table_name query_sql file_name(table_name 是表名，query_sql是 是否要手写querySql true or false,file_name是文件名)")
		}
		if i > 2 {
			log.Fatal("输入的参数最多只能为3个")
			continue
		}
		strs[i] = arg
	}
	ParseStrs(strs)
}

func ParseStrs(strs []string) {
	if len(strs) == 2 {
		parseBool, e := strconv.ParseBool(strs[1])
		if e != nil {
			log.Fatal("isTest 输入错误：", e)
		}
		var way code.Way = code.HToP{
			DataSource: code.Gamedb,
			TableName:  strs[0],
			IsTest:     parseBool,
		}
		if e = code.DataScrpitGenerate(way); e != nil {
			log.Fatal(e)
		}
	}

	if len(strs) == 3 {
		parseBool, e := strconv.ParseBool(strs[1])
		if e != nil {
			log.Fatal("query_sql 输入错误：", e)
		}
		var way code.Way = code.PToH{
			DataSource: code.Gamedb,
			TableName:  strs[0],
			QuerySql:   parseBool,
			FileName:   strs[2],
		}
		if e = code.DataScrpitGenerate(way); e != nil {
			log.Fatal(e)
		}
	}
}
