package auto_generate_model

type ModelGenerate struct {
	//数据库种类
	Dialect string
	//数据库需要参数
	DataSource string

	//表名
	TableName []string
	//创建文件夹路径
	Path string
	//需要创建的文件夹
	folders []string
}

func New(path string) *ModelGenerate {
	return &ModelGenerate{

	}
}

//AutoGenerate
//初始模板
//传入数据
// 1.读取数据库中内容
// 2.生成文件夹
// 3.生成文件
// 4.写入内容
func (a *ModelGenerate) AutoGenerate() error {
	// 1.读取数据库中内容
	a.LoadDataBase()

	return nil
}

func (a *ModelGenerate) LoadDataBase() {

}

// AutoFolder 自动生成文件夹
// 1.不需要自定义创建文件夹
func (a *ModelGenerate) AutoFolder() {

	//1.配置文件夹,同时生成文件夹
	defDir := []string{"${table_name}Controller", "${table_name}Service", "${table_name}Model", "${table_name}Utils"}
	a.folders = defDir

}
