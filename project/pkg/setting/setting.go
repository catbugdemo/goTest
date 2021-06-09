package setting

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Set struct {
	Redis Redis
	Postgre Postgre
}

type Redis struct {
	Host string `yaml:"host"`
	Password string`yaml:"password"`
	Timeout int `yaml:"timeout"`
	MaxActive int `yaml:"max_active"`
	MaxIdle int `yaml:"max_idle"`
	Db int
}

type Postgre struct {
	Type string `yaml:"type"`
	User string `yaml:"user"`
	Password string `yaml:"password"`
	Host string `yaml:"host"`
	Name string `yaml:"name"`
	TablePrefix string `yaml:"table_prefix"`
}

var Setting = Set{}

func InitSetting() {
	file, err := ioutil.ReadFile("project/conf/app.yml")
	if err != nil {
		log.Fatal("fail to read file:",err)
	}

	err = yaml.Unmarshal(file, &Setting)
	if err != nil {
		log.Fatal("fail to yaml unmarshal:",err)
	}

}
