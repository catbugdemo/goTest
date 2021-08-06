package main

import (
	"flag"
	"fmt"
	"log"
	"main.go/auto_project/data_scrpit/code"
	"strconv"
	"strings"
)

var params = flag.String("p", "", "string类型参数")

func main() {

	flag.Parse()

	s := *params

	split := strings.Split(s, "/")

	ParseStrs(split)
}

func ParseStrs(strs []string) {
	if len(strs) < 3 || len(strs) > 4 {
		fmt.Println("如果想生成 hdfs to pg 请输入3个参数如: ./main -p 数据库名/表名/true(isTest 选择测试服还是正式服)")
		fmt.Println("如果想生成 pg to hdfs 请输入4个参数如: ./main -p 数据库名/表名/true(querySql 选择是否自动生成)/filename(内部的一个参数数据)")
	}

	if len(strs) == 3 {
		db := code.Pointdb
		// 判断数据库db
		if strs[0] == "datadb" {
			db = code.Datadb
		}else{
			db["dbname"] = strs[0]
		}

		parseBool, e := strconv.ParseBool(strs[2])
		if e != nil {
			log.Fatal("isTest 输入错误：", e)
		}
		var way code.Way = code.HToP{
			DataSource: db,
			TableName:  strs[1],
			IsTest:     parseBool,
		}
		if e = code.DataScrpitGenerate(way); e != nil {
			log.Fatal(e)
		}
	}

	if len(strs) == 4 {
		db := code.Pointdb
		// 判断数据库db
		if strs[0] == "datadb" {
			db = code.Datadb
		}else{
			db["dbname"] = strs[0]
		}

		parseBool, e := strconv.ParseBool(strs[2])
		if e != nil {
			log.Fatal("query_sql 输入错误：", e)
		}
		var way code.Way = code.PToH{
			DataSource: db,
			TableName:  strs[1],
			QuerySql:   parseBool,
			FileName:   strs[3],
		}
		if e = code.DataScrpitGenerate(way); e != nil {
			log.Fatal(e)
		}
	}
}
