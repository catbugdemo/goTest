package code

import (
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
	"path"
	"strings"
)

var (
	//数据库连接参数
	Datadb = map[string]interface{}{
		"db":       "postgres",
		"host":     "49.234.137.226",
		"port":     "5432",
		"user":     "qipai",
		"dbname":   "datadb",
		"sslmode":  "disable",
		"password": "qipai#xq5",
	}

	Pointdb = map[string]interface{}{
		"db":       "postgres",
		"host":     "123.206.176.76",
		"port":     "5432",
		"user":     "qipai",
		"dbname":   "",
		"sslmode":  "disable",
		"password": "qipai#xq5",
	}
)

type HToP struct {
	// 数据库连接参数
	DataSource map[string]interface{}
	// 是否是测试服 ,不是则为正式服
	//	测试："hdfs://bigdata1:8020"
	//	正式："hdfs://zonst-bigdata"
	IsTest bool
	// 数据库表名
	TableName string
	// 数据库名称
	DbName string
	// 自定义表名
	SelfTableName string
}

type PToH struct {
	// 数据库数据
	DataSource map[string]interface{}
	// 是否生成 query_sql
	QuerySql bool
	// 表名
	TableName string
	// 拼接参数
	FileName string
	// 是否是测试服 ,不是则为正式服
	//	测试："hdfs://bigdata1:8020"
	//	正式："hdfs://zonst-bigdata"
	IsTest bool
}

type Way interface {
	// 初始化数据
	initStruct()
	// 创文件夹，文件
	create() error
	// 格式化输出string
	formatJSON() (string, error)
	// 谢润文件
	writeFile(buf []byte) error
}

// 入参接口
// 1.创建文件夹,文件
// 2.生成格式化代码
// 3.写入文件
func DataScrpitGenerate(way Way) error {
	// 创建文件文件夹
	if e := way.create(); e != nil {
		return errorx.Wrap(e)
	}
	// 生成格式化代码
	json, e := way.formatJSON()
	if e != nil {
		return errorx.Wrap(e)
	}
	// 写入文件
	if e = way.writeFile([]byte(json)); e != nil {
		return errorx.Wrap(e)
	}

	return nil
}

// 数据初始化
func (h HToP) initStruct() {
	if h.DataSource == nil {
		h.DataSource = map[string]interface{}{}
	}
	initDataSource(h.DataSource)

}

func initDataSource(db map[string]interface{}) {
	if _, ok := db["db"]; !ok {
		db["db"] = "postgres"
	}
	if _, ok := db["host"]; !ok {
		db["host"] = "127.0.0.1"
	}
	if _, ok := db["port"]; !ok {
		db["port"] = "5432"
	}
	if _, ok := db["user"]; !ok {
		db["user"] = "postgres"
	}
	if _, ok := db["dbname"]; !ok {
		db["dbname"] = "db_name"
	}
	if _, ok := db["sslmode"]; !ok {
		db["sslmode"] = "disable"
	}
	if _, ok := db["password"]; !ok {
		db["password"] = ""
	}
}

// hdfs to pg 文件生成
func (h HToP) create() error {
	// 获取本地路径
	getwd, e := os.Getwd()
	if e != nil {
		return errorx.Wrap(e)
	}
	dir := path.Join(getwd, "auto_scrpit", "hdfs_to_pg")
	// 生成文件夹
	if e = os.MkdirAll(dir, os.ModePerm); e != nil {
		return errorx.Wrap(e)
	}
	// 生成文件
	if _, e = os.Create(path.Join(dir, "hdfs_to_pg.json")); e != nil {
		return errorx.Wrap(e)
	}
	return nil
}

// 根据数据库获取相应数据
// 格式化输出
// ${index_and_type} index 和 type 类型，string 特殊处理
// ${test_path} 是否是测试服地址
// ${path} 地址转换
// ${columns} 参数匹配 column
// ${username} 用户名
// ${password} 密码
// ${db_path}
// ${table_name}
func (h HToP) formatJSON() (string, error) {
	columns, e := getDatabase(h.DataSource, h.TableName)
	if e != nil {
		return "", errorx.Wrap(e)
	}
	var str = `
--- hdfs 数据到 pg

{
    "job": {
        "content": [
            {
                "reader": {
                    "name": "hdfsreader",
                    "parameter": {
                        "column": [
${index_and_type}
                        ],
                        //正式环境确定地址后不再变化，测试环境固定hdfs://bigdata1:8020
                        "defaultFS": "${test_path}",
                        "encoding": "UTF-8",
                        "fieldDelimiter": "\u0001",
                        "fileType": "text",
                        //每张表地址不同
                        "path": "${path}"
                    }
                },
                "writer": {
                    "name": "postgresqlwriter",
                    "parameter": {
                        "column": [
${columns}
                        ],
                        "connection": [
                            {
                                "jdbcUrl": "jdbc:postgresql://${db_path}",
                                "table": [                             
                                    "${table_name}"
                                ]
                            }
                        ],
                        "password": "${password}",
                        "postSql": [],
                        "preSql": [],
                        "username": "${username}"
                    }
                }
            }
        ],
        "setting": {
            "speed": {
                "channel": "1"
            }
        }
    }
}
`
	str = strings.ReplaceAll(str, "${index_and_type}", indexAndType(columns))
	str = strings.ReplaceAll(str, "${test_path}", testPath(h.IsTest))
	str = strings.ReplaceAll(str, "${path}", fmt.Sprintf("/user/hive/warehouse/%s.db/%s/dt=${dt}", h.DbName, h.SelfTableName))
	str = strings.ReplaceAll(str, "${columns}", returnColumns(columns))
	str = strings.ReplaceAll(str, "${username}", h.DataSource["user"].(string))
	str = strings.ReplaceAll(str, "${password}", h.DataSource["password"].(string))
	str = strings.ReplaceAll(str, "${db_path}",
		fmt.Sprintf("%s:%s/%s", h.DataSource["host"], h.DataSource["port"], h.TableName))
	str = strings.ReplaceAll(str, "${table_name}", h.TableName)
	return str, nil
}

func (h HToP) writeFile(buf []byte) error {
	getwd, e := os.Getwd()
	if e != nil {
		return errorx.Wrap(e)
	}
	file, e := os.OpenFile(path.Join(getwd, "auto_scrpit", "hdfs_to_pg", "hdfs_to_pg.json"), os.O_CREATE|os.O_APPEND|os.O_RDWR, 6)
	defer file.Close()
	if e != nil {
		return errorx.Wrap(e)
	}
	if _, e = file.Write(buf); e != nil {
		return errorx.Wrap(e)
	}
	return nil
}

func indexAndType(column []Column) string {
	tmp := ""
	for i := 0; i < len(column); i++ {
		tmp += fmt.Sprintf("\t\t\t\t\t\t\t{\n\t\t\t\t\t\t\t\t\"index\":%d,\n\t\t\t\t\t\t\t\t\"type\":\"%s\"\n\t\t\t\t\t\t\t},\n",
			i, returnType(column[i].ColumnType))
	}
	return tmp[:len(tmp)-1]
}

func testPath(isTest bool) string {
	if isTest {
		return "hdfs://bigdata1:8020"
	}
	return "hdfs://zonst-bigdata"
}

func returnColumns(columns []Column) string {
	tmp := ""
	for _, column := range columns {
		tmp += fmt.Sprintf("\t\t\t\t\t\t\"%s\",\n", column.ColumnName)
	}
	return tmp[:len(tmp)-2]
}

func returnType(columnType string) string {
	switch columnType {
	case "integer", "serial":
		return "int"
	case "double precision", "money", "numeric", "real":
		return "double"
	case "varchar", "char", "text", "bit", "inet":
		return "string"
	case "date", "time", "timestamp", "timestamp with time zone":
		return "date"
	case "bool":
		return "boolean"
	case "bytea":
		return "bytes"
	default:
		return columnType
	}
}

func (p PToH) initStruct() {
	if p.DataSource == nil {
		p.DataSource = map[string]interface{}{}
	}
	initDataSource(p.DataSource)
}

func (p PToH) create() error {
	// 获取本地路径
	getwd, e := os.Getwd()
	if e != nil {
		return errorx.Wrap(e)
	}
	dir := path.Join(getwd, "auto_scrpit", "pg_to_hdfs")
	// 生成文件夹
	if e = os.MkdirAll(dir, os.ModePerm); e != nil {
		return errorx.Wrap(e)
	}
	// 生成文件
	if _, e = os.Create(path.Join(dir, "pg_to_hdfs.json")); e != nil {
		return errorx.Wrap(e)
	}
	return nil
}

// ${username} 数据库用户
// ${password} 密码
// ${pg_params}
// ${query_sql} 内容
// ${table_name}表名
// ${pg_column}
// ${file_name}
// ${is_test}
func (p PToH) formatJSON() (string, error) {
	columns, e := getDatabase(p.DataSource, p.TableName)
	if e != nil {
		return "", errorx.Wrap(e)
	}
	var str = `
--- pg 到 hdfs

{
    "job": {
        "content": [
            {
                "reader": {
                     "name": "postgresqlreader",
                    "parameter": {
                        "username": "${username}",
                        "password": "${password}",
                        "connection": [
                            {
                                "querySql": [
                                     "select 
${query_sql}   
                                   from ${table_name} where reg_date='${day}'"
                                ],
                                "jdbcUrl": [
                                    "jdbc:postgresql://${db_path}"
                                ]
                            }
                        ]
                    }
                },
                "writer": {
                    "name": "hdfswriter",
                    "parameter": {
                        "column": [
${pg_column}
                        	],
                        "compress": "NONE",
                        "defaultFS": "${is_test}",
                        "fieldDelimiter": "\u0001",
                        "fileName": "${day}_${flie_name}",
                        "fileType": "orc",
                        "path": "${path}",
                        "writeMode": "append"
                    }
                }
            }
        ],
        "setting": {
            "speed": {
                "channel": "1"
            }
        }
    }
}
`
	str = strings.ReplaceAll(str, "${username}", p.DataSource["user"].(string))
	str = strings.ReplaceAll(str, "${password}", p.DataSource["password"].(string))
	str = strings.ReplaceAll(str, "${table_name}", p.TableName)
	str = strings.ReplaceAll(str, "${path}", fmt.Sprintf("/user/hive/warehouse/%s.db/%s/dt=${dt}", p.DataSource["dbname"], p.TableName))
	str = strings.ReplaceAll(str, "${db_path}",
		fmt.Sprintf("%s:%s/%s", p.DataSource["host"], p.DataSource["port"], p.TableName))
	tmp := ""
	if p.QuerySql {
		tmp = returnQuerySql(columns)
	}
	str = strings.ReplaceAll(str, "${query_sql}", tmp)
	str = strings.ReplaceAll(str, "${pg_column}", returnPGColumn(columns))
	str = strings.ReplaceAll(str, "${flie_name}", p.FileName)

	//	测试："hdfs://bigdata1:8020"
	//	正式："hdfs://zonst-bigdata"
	tmp = "hdfs://zonst-bigdata"
	if p.IsTest {
		tmp = "hdfs://bigdata1:8020"
	}
	str = strings.ReplaceAll(str,"${is_test}",tmp)

	return str, nil
}

func (p PToH) writeFile(buf []byte) error {
	getwd, e := os.Getwd()
	if e != nil {
		return errorx.Wrap(e)
	}
	file, e := os.OpenFile(path.Join(getwd, "auto_scrpit", "pg_to_hdfs", "pg_to_hdfs.json"), os.O_CREATE|os.O_APPEND|os.O_RDWR, 6)
	defer file.Close()
	if e != nil {
		return errorx.Wrap(e)
	}
	if _, e = file.Write(buf); e != nil {
		return errorx.Wrap(e)
	}
	return nil
}

func returnQuerySql(columns []Column) string {
	tmp := ""
	for _, column := range columns {
		tmp += fmt.Sprintf("\t\t\t\t\t\t\t\t\t\t%s,\n", column.ColumnName)
	}
	return tmp[:len(tmp)-2]
}

func returnPGColumn(columns []Column) string {
	tmp := ""
	for _, column := range columns {
		tmp += fmt.Sprintf("\t\t\t\t\t\t\t\t{\n\t\t\t\t\t\t\t\t\t\"name\":\"%s\",\n\t\t\t\t\t\t\t\t\t\"type\":\"%s\"\n\t\t\t\t\t\t\t\t},\n",
			column.ColumnName, returnPGType(column))
	}
	return tmp[:len(tmp)-2]
}

func returnPGType(column Column) string {
	switch column.ColumnType {
	case "integer", "smallint", "tinyint", "int", "smallserial":
		return "int"
	case "bigint", "serial", "bigserial", "money":
		return "bigint"
	case "decimal", "numeric", "real":
		return "float"
	case "bool":
		return "boolean"
	case "character", "text", "json", "jsonb", "array":
		return "string"
	case "bytes":
		return "bytea"
	default:
		return column.ColumnType
	}
}

// 数据库列属性
type Column struct {
	ColumnName string `gorm:"column:column_name"` // column_name
	ColumnType string `gorm:"column:column_type"` // column_type
}

func getDatabase(db map[string]interface{}, tableName string) ([]Column, error) {
	switch db["db"] {
	case "postgres", "pg", "psql":
		columns, e := getPostgresqlColumn(db, tableName)
		if e != nil {
			return nil, errorx.Wrap(e)
		}
		return columns, nil
	case "mysql":
	default:
		return nil, errorx.NewFromString("db 输入错误:postgres 或者 mysql")
	}
	return nil, nil
}

func getPostgresqlColumn(db map[string]interface{}, tableName string) ([]Column, error) {
	var FindColumnsSql = `
        SELECT
            a.attnum AS column_number,
            a.attname AS column_name,
            --format_type(a.atttypid, a.atttypmod) AS column_type,
            a.attnotnull AS not_null,
			COALESCE(pg_get_expr(ad.adbin, ad.adrelid), '') AS default_value,
    		COALESCE(ct.contype = 'p', false) AS  is_primary_key,
    		CASE
        	WHEN a.atttypid = ANY ('{int,int8,int2}'::regtype[])
          		AND EXISTS (
				SELECT 1 FROM pg_attrdef ad
             	WHERE  ad.adrelid = a.attrelid
             	AND    ad.adnum   = a.attnum
             	-- AND    ad.adsrc = 'nextval('''
                --	|| (pg_get_serial_sequence (a.attrelid::regclass::text
                --	                          , a.attname))::regclass
                --	|| '''::regclass)'
             	)
            THEN CASE a.atttypid
                    WHEN 'int'::regtype  THEN 'serial'
                    WHEN 'int8'::regtype THEN 'bigserial'
                    WHEN 'int2'::regtype THEN 'smallserial'
                 END
			WHEN a.atttypid = ANY ('{uuid}'::regtype[]) AND COALESCE(pg_get_expr(ad.adbin, ad.adrelid), '') != ''
            THEN 'autogenuuid'
        	ELSE format_type(a.atttypid, a.atttypmod)
    		END AS column_type
		FROM pg_attribute a
		JOIN ONLY pg_class c ON c.oid = a.attrelid
		JOIN ONLY pg_namespace n ON n.oid = c.relnamespace
		LEFT JOIN pg_constraint ct ON ct.conrelid = c.oid
		AND a.attnum = ANY(ct.conkey) AND ct.contype = 'p'
		LEFT JOIN pg_attrdef ad ON ad.adrelid = c.oid AND ad.adnum = a.attnum
		WHERE a.attisdropped = false
		AND n.nspname = 'public'
		AND c.relname = ?
		AND a.attnum > 0
		ORDER BY a.attnum
	`
	var s = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s", db["host"], db["port"], db["user"], db["dbname"], db["sslmode"], db["password"])
	open, e := gorm.Open("postgres", s)
	if e != nil {
		return nil, errorx.Wrap(e)
	}
	open.SingularTable(true)
	var columns = make([]Column, 0, 10)
	open.Raw(FindColumnsSql, tableName).Find(&columns)
	if len(columns) < 1 {
		return nil, errorx.NewFromString("表名 或者数据库参数存在问题")
	}
	return columns, nil
}
