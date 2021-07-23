package code

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	//数据库连接参数
	gamedb = map[string]interface{}{
		"db":       "postgres",
		"host":     "49.234.137.226",
		"port":     "5432",
		"user":     "zonst_xyx",
		"dbname":   "game_test",
		"sslmode":  "disable",
		"password": "zonst_xyx_fengtao",
	}
)

func TestDataScriptGenerate(t *testing.T) {
	// hdfs_to_pg
	var h = HToP{
		// 数据库参数
		DataSource: gamedb,
		// 是否是测试服地址
		IsTest: true,
		// 数据库表名
		TableName: "user_send_process",
	}

	// pg_to_hdfs
	var p = PToH{
		// 数据库参数
		DataSource: gamedb,
		// 是否自动生成 QuerySql
		QuerySql:   true,
		// 表名
		TableName:  "user_send_process",
		// 文件名称
		FileName:   "123",
	}
	// 会在该文件目录下自动生成一个 auto_scrpit 包
	t.Run("hdfs_to_pg", func(t *testing.T) {
		var way Way = h
		e := DataScrpitGenerate(way)
		assert.Nil(t, e)
	})
	t.Run("pg_to_hdfs", func(t *testing.T) {
		var way Way = p
		e := DataScrpitGenerate(way)
		assert.Nil(t, e)
	})
}