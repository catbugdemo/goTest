package setting

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Set struct {
	Redis Redis
}

type Redis struct {
	Host string
	Password string
	Timeout int
	MaxActive int
	MaxIdle int
	Db int
}

var Setting = Set{}

func InitSetting() {
	file, err := ioutil.ReadFile("../conf/app.yml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(file, &Setting)
	if err != nil {
		panic(err)
	}
}
