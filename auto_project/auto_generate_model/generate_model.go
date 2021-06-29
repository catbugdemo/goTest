package auto_generate_model

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
)

type ReplaceType struct {
	Path   string
	ReqSrc interface{}
	Src    interface{}

	//是否要自己写sql
	WriteMyself bool
	//数据库结构体名称
	NameSrc string
	//请求结构体名称
	NameReqSrc string
	//文件夹名称
	NameDir string
	//自定义表名
	NameTable string
}

func New(path string, reqSrc interface{}, src interface{}, writeMyself bool) (*ReplaceType, error) {
	if reqSrc == nil || src == nil {
		return nil, errors.New("reqSrc or src is nil")
	}
	//数据库结构体名称
	srcTypeOf := reflect.TypeOf(src)
	nameSrc := srcTypeOf.Name()
	//请求结构体名称
	reqSrcTypeOf := reflect.TypeOf(reqSrc)
	nameReqSrc := reqSrcTypeOf.Name()
	//格式化文件夹名称
	nameDir := strings.ToLower(string(srcTypeOf.Name()[0])) + srcTypeOf.Name()[1:]
	//自定义表名
	nameTable := addAndTransferToLower(nameDir)

	return &ReplaceType{
		Path:        path,
		ReqSrc:      reqSrc,
		Src:         src,
		WriteMyself: writeMyself,
		NameSrc:     nameSrc,
		NameReqSrc:  nameReqSrc,
		NameDir:     nameDir,
		NameTable:   nameTable,
	}, nil
}

// GenerateData 自动创建大数据类型 (path：绝对参数,reqSrc:请求结构体，src:数据库结构体)
// 1.创建文件夹,文件
// 2.编写内容
// 3.将内容写入文件
func GenerateData(path string, reqSrc interface{}, src interface{}, writeMyself bool) error {
	replaceType, e := New(path, reqSrc, src, writeMyself)
	if e != nil {
		return e
	}

	//创建文件夹
	if e = replaceType.create(); e != nil {
		return e
	}

	//写入内容
	if e = replaceType.write(); e != nil {
		return e
	}

	return nil
}

// create
// 创建文件夹，创建文件
func (r ReplaceType) create() error {
	//对名称进行格式化
	pathDir := r.Path + "/" + r.NameDir

	//创建文件夹
	if e := os.MkdirAll(pathDir, os.ModePerm); e != nil {
		log.Println("mkdir failed")
		return e
	}

	//创建文件
	creFile := []string{"router.go", "service.go", "model.go"}
	for _, v := range creFile {
		_, e := os.Create(pathDir + "/" + v)
		if e != nil {
			return e
		}

	}
	return nil
}

// write 写入内容
//
func (r ReplaceType) write() error {
	if e := r.writeModel(); e != nil {
		return e
	}
	if e := r.writeService(); e != nil {
		return e
	}
	if e := r.writeRouter(); e != nil {
		return e
	}
	return nil
}

// writeModel 写入model
// ${package_name} 包名
// ${struct} 结构体
// ${module_req} 根据请求结构体自动生成
// ${req_struct_name} 请求结构体名称
// ${struct_name} 结构体名称
// ${table_name} 表名
// ${middle} 中间改变
func (r ReplaceType) writeModel() error {

	var front = `package ${name_dir}

import "datasrv/dependence/db"

//根据请求结构体生成
${module_req}

// 数据库表结构体
${struct}

`
	auto := `	
// 自动生成代码结构
func Get${struct_name}QueryList(model ${req_struct_name}) ([]${struct_name}, error) {
	a := make([]${struct_name}, 0)
	if err := db.DataDB.Table("${table_name}").
${middle}
		return a, err
	}
	return a, nil
}
`
	external := `
// 开发人员自己写sql
func Get${struct_name}QueryList(model ${struct_name}) ([]${struct_name}, error) {
	a := make([]${struct_name}, 0)
	if err := db.DataDB.Table("${table_name}").Raw("").Find(&a).Error; err != nil {
		return a, err
	}
	return a, nil
}
`
	if r.WriteMyself {
		front += external
	} else {
		front += auto
	}

	front = strings.ReplaceAll(front, "${name_dir}", r.NameDir)
	front = strings.ReplaceAll(front, "${struct_name}", r.NameSrc)
	front = strings.ReplaceAll(front, "${req_struct_name}", r.NameReqSrc)
	front = strings.ReplaceAll(front, "${struct}", formatPrint(r.Src))
	front = strings.ReplaceAll(front, "${middle}", middleGenerate(r.ReqSrc))
	front = strings.ReplaceAll(front, "${table_name}", r.NameTable)

	tmp := formatPrint(r.ReqSrc)
	tmp = strings.ReplaceAll(tmp, r.NameReqSrc, r.NameSrc+"ModuleReq")
	front = strings.ReplaceAll(front,"${module_req}",tmp)

	if e := openAndWriteStringFile(r.Path+"/"+r.NameDir+"/model.go", front); e != nil {
		return e
	}
	return nil
}

func openAndWriteStringFile(pathFile string, write string) error {
	file, e := os.OpenFile(pathFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 6)
	defer file.Close()
	if e != nil {
		return e
	}
	if _, e = file.WriteString(write); e != nil {
		return e
	}
	return nil
}

// writeService 写入service
// ${name_dir} 包名
// ${req_struct} 请求结构体
// ${req_struct_name} 请求结构体名称
// ${struct_name} 结构体名称
// ${transfer_param} 需要传递的参数
func (r ReplaceType) writeService() error {
	var front = `package ${name_dir}

import (
    "datasrv/utils/logstat"
    "github.com/gin-gonic/gin"
    "net/http"
    "time"
    "zonst/logging"
)



func ${struct_name}QueryList(c *gin.Context) {

    defer logstat.LogStat("${struct_name}QueryList", c, time.Now())
  	// 请求结构体
	${req_struct}
  
    req := ${req_struct_name}{}
  
    if err := c.ShouldBindJSON(&req); err != nil {
        logging.Errorf("${req_struct_name},err:%v\n", err.Error())
        c.JSON(http.StatusOK, gin.H{"error": -1, "message": err.Error()})
        return
    }
  
  // model 层函数 get + 当前方法名
  
    data, err := Get${struct_name}QueryList(${struct_name}ModuleReq{
${transfer_param}
	})
  
    if err != nil {
        logging.Errorf("Get${struct_name}List,err:%v\n", err.Error())
        c.JSON(http.StatusOK, gin.H{"error": -1, "message": err.Error()})
        return
    }
  
    c.JSON(http.StatusOK, gin.H{"error": 0, "message": "", "data": data})
}
`
	front = strings.ReplaceAll(front, "${name_dir}", r.NameDir)
	front = strings.ReplaceAll(front, "${req_struct}", formatPrint(r.ReqSrc))
	front = strings.ReplaceAll(front, "${req_struct_name}", r.NameReqSrc)
	front = strings.ReplaceAll(front, "${struct_name}", r.NameSrc)
	front = strings.ReplaceAll(front, "${transfer_param}", transferParam(r.ReqSrc))

	if e := openAndWriteStringFile(r.Path+"/"+r.NameDir+"/service.go", front); e != nil {
		return e
	}
	return nil
}

func transferParam(src interface{}) string {
	typeOf := reflect.TypeOf(src)
	tmp := ""
	for i := 0; i < typeOf.NumField(); i++ {
		tmp += fmt.Sprintf("\t\t%s:    req.%s,\n", typeOf.Field(i).Name, typeOf.Field(i).Name)
	}
	return tmp
}

// writeRouter 路由编写
// ${name_dir} 包名
// ${struct_name} 结构体名称
// ${path_route} 路由路径
func (r ReplaceType) writeRouter() error {
	var front = `package ${name_dir}

import "github.com/gin-gonic/gin"

func Router(r gin.IRouter) {
	r.POST("${path_route}", ${struct_name}QueryList)   
}
`
	split := strings.Split(addAndTransferToLower(r.NameSrc), "_")
	tmp := "/"
	for _, v := range split {
		tmp += v + "/"
	}
	tmp += "query-list"

	front = strings.ReplaceAll(front, "${name_dir}", r.NameDir)
	front = strings.ReplaceAll(front, "${struct_name}", r.NameSrc)
	front = strings.ReplaceAll(front, "${path_route}", tmp)

	if e := openAndWriteStringFile(r.Path+"/"+r.NameDir+"/router.go", front); e != nil {
		return e
	}
	return nil
}

// formatPrint 格式化输出结构体
// 问题：未结构化生成结构体
func formatPrint(src interface{}) string {
	typeOf := reflect.TypeOf(src)
	tmp := ""
	for i := 0; i < typeOf.NumField(); i++ {
		tmp += fmt.Sprintf("    %s  %s    `%s` \n", typeOf.Field(i).Name, typeOf.Field(i).Type, typeOf.Field(i).Tag)
	}
	return fmt.Sprintf("type %s struct {\n%s \n}", typeOf.Name(), tmp)
}

// addAndTransferToLower 现将第一个字符转换为小写，大写前添加下划线并全部转为小写,
func addAndTransferToLower(name string) string {
	name = strings.ToLower(string(name[0])) + name[1:]
	index := 0
	for k, v := range name {
		if v >= 65 && v <= 90 {
			name = name[0:k+index] + "_" + name[k+index:]
			index++
		}
	}
	return strings.ToLower(name)
}

// 创建model中的Where语句
func middleGenerate(src interface{}) string {
	typeOf := reflect.TypeOf(src)
	tmp := ""
	for i := 0; i < typeOf.NumField(); i++ {
		if typeOf.Field(i).Name == "CalDate" {
			tmp += fmt.Sprintf("Where(\"cal_date between ? and ?\", model.StartTime, model.EndTime).\n")
		}
		tmp += fmt.Sprintf("\t\tWhere(\"%s=?\", model.%s).", addAndTransferToLower(typeOf.Field(i).Name), typeOf.Field(i).Name)
		if i == typeOf.NumField()-1 {
			tmp += fmt.Sprintf("Find(&a).Error; err != nil {")
		}
		tmp += fmt.Sprintf("\n")
	}
	return tmp
}
