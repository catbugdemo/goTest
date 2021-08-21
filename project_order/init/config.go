package init

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"log"
	"sync"
)

var (
	config *viper.Viper
	mu     sync.RWMutex
)

type Config struct {

}

func init() {
	InitConfig()
}

func InitConfig() {
	v := viper.New()
	// 初始化配置信息
	v.SetConfigName("local")
	v.SetConfigType("toml")
	v.AddConfigPath("$GOPATH/src/github.com/catbugdemo/goTest/")
	if err := ReadInConfig(v); err != nil {
		log.Fatalf("读取配置失败：%v\n", err)
	}

	// 当配置改变是重新运行
	v.WatchConfig()
	/*	v.OnConfigChange(func(in interface{}) {

		})*/

	config = v
}

// ReadInConfig 通过读写锁匹配内容
func ReadInConfig(v *viper.Viper) error {
	mu.RLock()
	defer mu.RUnlock()
	if err := v.ReadInConfig(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// GetConfig 获取配置
func GetConfig() *viper.Viper {
	mu.RLock()
	defer mu.RUnlock()
	return config
}
