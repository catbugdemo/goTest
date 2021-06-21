package auto

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"
)

type Field struct {
	//是否为空
	isZero bool
	//种类名称
	TypeName string
	//字段
	FieldName string
	//标签名
	TagName string
	//值
	Value interface{}
}

func GenerateListWhere(src interface{}, withListArgs bool, replacement ...map[string]string) string {
	if len(replacement) == 0 {
		replacement = []map[string]string{
			map[string]string{},
		}
	}
	//var format = `
   //${handle}
//`
	vValue := reflect.ValueOf(src)
	vType := reflect.TypeOf(src)
	f := func() {
		for i := 0; i < vValue.NumField(); i++ {
			// 判断结构体数据是否属于传统类型
			// only handle basic types, otherwise continue
			if !in(vType.Field(i).Type.String(), []string{
				"int", "int8", "int16", "int32", "int64",
				"float32", "float64",
				"string",
				"uint8", "uint16", "uint32", "uint64",
				"time.Time",
			}) {
				continue
			}

			//不属于传统值类型
			//是否是time.Time
			if vType.Field(i).Type.AssignableTo(reflect.TypeOf(time.Time{})) {

			}
		}
	}
	fmt.Println(f)
	return ""
}

//GenerateAddOneAPI
//转入数据:结构体src.HTTPAdd+模块名
func GenerateAddOneAPI(src interface{}, replacement ...map[string]string) string {
	if len(replacement) == 0 {
		replacement = []map[string]string{
			map[string]string{},
		}
	}

	//将结构体 src 和 函数名传入
	handleDefault(src, replacement[0])

	//自定义HTTPAdd
	var resultf = `
// Auto generate by github.com/fwhezfwhez/model_convert.GenerateAddOneAPI().
func ${handler_name} (c *gin.Context) {
    var param ${model}
    if e := c.Bind(&param); e!=nil {
        c.JSON(400, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }

    if e:=(${db_instance}.Model(&${model}{}).Create(&param).Error); e!=nil {
        ${handle_error}
        c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }

    if param.RedisKey() != "" {
        conn := ${redis_conn}
        defer conn.Close()
        param.DeleteFromRedis(conn)
    }
    c.JSON(200, gin.H{"message": "success", "data": param})
}
`
	//将结果中的数据进行替换
	result := strings.Replace(resultf, "${handler_name}", replacement[0]["${handler_name}"], -1)
	result = strings.Replace(result, "${handle_error}", replacement[0]["${handle_error}"], -1)
	result = strings.Replace(result, "${db_instance}", replacement[0]["${db_instance}"], -1)
	result = strings.Replace(result, "${model}", replacement[0]["${model}"], -1)
	result = strings.Replace(result, "${redis_conn}", replacement[0]["${redis_conn}"], -1)
	return result
}

//GenerateListAPI
// Replacement optional as:
// - ${page} "page"
// - ${size} "size"
// - ${order_by} ""
// - ${util_pkg} "util"
// - ${db_instance} "db.DB"
// - ${handler_name} "HTTPListUser"
// - ${model} "model.User"
// - ${handle_error} "fmt.Println(e)"
// - ${jump_fields}, "password,pw"
// - ${layout}, "2006-01-02 15:04:03"
// - ${time_zone} "time.Local"
// 传入数据
func GenerateListAPI(src interface{}, withListArgs bool, replacement ...map[string]string) string {
	if len(replacement) == 0 {
		replacement = []map[string]string{
			map[string]string{},
		}
	}

	if replacement[0]["${page}"] == "" {
		replacement[0]["${page}"] = "page"
	}

	if replacement[0]["${size}"] == "" {
		replacement[0]["${size}"] = "size"
	}

	if replacement[0]["${db_instance}"] == "" {
		replacement[0]["${db_instance}"] = "db.DB"
	}

	if replacement[0]["${util_pkg}"] == "" {
		replacement[0]["${util_pkg}"] = "mc_util"
	}
	if replacement[0]["${list_layout}"] == "" {
		replacement[0]["${list_layout}"] = "2006-01-02"
	}

	if replacement[0]["${handler_name}"] == "" {
		replacement[0]["${handler_name}"] = "HTTPListUser"
	}
	if replacement[0]["${list_layout}"] == "" {
		replacement[0]["${list_layout}"] = "2006-01-02"
	}

	//模块名
	vType := reflect.TypeOf(src)
	if replacement[0]["${model}"] == "" {
		replacement[0]["${model}"] = vType.String()
	}

	if replacement[0]["${handle_error}"] == "" {
		replacement[0]["${handle_error}"] = "log.Println(e)"
	}

	var copyMap = make(map[string]string)
	for k, v := range replacement[0] {
		copyMap[k] = v
	}
	return ""
}

// GenerateCRUD 自动打印CRUD
// 传入的内容：需要的结构体，外部的一些导入类
func GenerateCRUD(src interface{}, replacement ...map[string]string) string {
	//判断是否有外部数据
	if len(replacement) == 0 {
		//创建数据，replacement
		replacement = []map[string]string{
			map[string]string{},
		}
	}

	//获取结构体中所有的种类
	vType := reflect.TypeOf(src)
	//获取结构体中所有的值
	vValue := reflect.ValueOf(src)

	//该循环用来处理标签字符串json中的匹配字段
	//依次处理结构体中所有的种类
	for i := 0; i < vType.NumField(); i++ {
		//Get方法返回标签字符串中键key对应的值。如果标签中没有该键，会返回"",如果标签不符合标准格式，Get的返回值是不确定的。
		tagStr := vType.Field(i).Tag.Get("json")
		//如果不存在则跳过该循环
		if tagStr == "-" || tagStr == "" {
			continue
		}
		//取第一个json来进行匹配
		arr := strings.Split(tagStr, ",")
		//返回将s前后端所有空白（unicode.IsSpace指定）都去掉的字符串
		tagValue := strings.TrimSpace(arr[0])
		//返回当前持有的值
		valueI := vValue.Field(i).Interface()

		var rangement = make([]Field, 0, 10)
		rangement = append(rangement, Field{
			//根据值的种类判断,
			isZero: IfZero(valueI),
			//当前属性的类型名 id int 中的int
			TypeName: vType.Field(i).Type.Name(),
			//当前属性名 id int 中的 id
			FieldName: vType.Field(i).Name,
			// json字段中的与前端对应的名称
			TagName: tagValue,
			//获取结构体中的值
			Value: valueI,
		})
	}

	//模块名称
	modelName := vType.Name()
	//定义add模块名
	replacement[0]["${handler_name}"] = "HTTPAdd" + modelName
	GenerateAddOneAPI(src, replacement...)

	//定义获list的值
	replacement[0]["${handler_name}"] = "HTTPList" + modelName
	GenerateListAPI(src, false, replacement...)
	return ""
}

//IfZero 判断结构体重的数据是否为空
//传入某一结构体自定义的值
func IfZero(arg interface{}) bool {
	//判断是否为空
	if arg == nil {
		return true
	}

	//获取数据的数据种类
	switch v := arg.(type) {
	case int, int16, int32, int64:
		if v == 0 {
			return true
		}
	case float32:
		//将float32转为float64，并进行绝对值的判断
		r := float64(v)
		return math.Abs(r-0) < 0.0000001
	case float64:
		return math.Abs(v-0) < 0.0000001
		//如果是指针类型的数据
	case *string, *int, *int64, *int32, *int16, *int8, *float32, *float64, *time.Time:
		if v == nil {
			return true
		}
		//如果是时间如何判断，未完成
	case time.Time:
		return false
	default:
		return false
	}
	return false
}

//handleDefault 默认处理
//传入数据：结构体src,replacement
func handleDefault(src interface{}, replacement map[string]string) {
	//制定页码
	if replacement["${page}"] == "" {
		replacement["${page}"] = "page"
	}

	//制定页码中的尺寸
	if replacement["${size}"] == "" {
		replacement["${size}"] = "size"
	}

	//制定数据库实例
	if replacement["${db_instance}"] == "" {
		replacement["${db_instance}"] = "db.DB"
	}

	//制定工具包
	if replacement["${util_pkg}"] == "" {
		replacement["${util_pkg}"] = "util"
	}

	//制定
	if replacement["${handler_name}"] == "" {
		replacement["${handler_name}"] = "HTTPListUser"
	}

	//制定模块中的数据
	vType := reflect.TypeOf(src)
	//确定模块名 (包名+结构体名): auto.T
	if replacement["${model}"] == "" {
		replacement["${model}"] = vType.String()
	}
	arr := strings.Split(replacement["${model}"], ".")
	if len(arr) == 2 {
		//包名
		replacement["${pkg_name_prefix}"] = arr[0] + "."
		//结构体名
		replacement["${struct_name}"] = arr[1]
	} else if len(arr) == 1 {
		replacement["${struct_name}"] = arr[0]
	} else {
		//do nothing
	}

	//解析到包名
	if replacement["${generate_to_pkg}"] == "" {
		replacement["${generate_to_pkg}"] = strings.ToLower(replacement["struct_name"])
	}
	//异常处理
	if replacement["${handle_error}"] == "" {
		replacement["${handle_error}"] = "fmt.Println(e, string(debug.Stack()))"
	}
	//缓存连接工具
	if replacement["${redis_conn}"] == "" {
		replacement["${redis_conn}"] = "redistool.RedisPool.Get()"
	}
}

func in(s string, arr []string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}
