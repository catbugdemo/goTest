package auto_generate_model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// 请求结构体

type AppCoinTaxReq struct {
	StartTime    string `json:"start_date"`
	EndTime      string `json:"end_date"`
	GameId       int    `json:"game_id"`
	GameAreaId   int    `json:"game_area_id"`
	AppChannel   string `json:"app_channel"`
	LoginChannel string `json:"login_channel"`
}
type AppCoinTax struct {
	CalDate        string `json:"cal_date" gorm:"cal_date"`
	GameId         int    `json:"game_id" gorm:"game_id"`                   // 用户平台id
	GameAreaId     int    `json:"game_area_id" gorm:"game_area_id"`         // --子游戏id
	AppChannel     string `json:"app_channel" gorm:"app_channel"`           // --子渠道
	LoginChannel   string `json:"login_channel" gorm:"login_channel"`       // --登陆渠道
	GameRule       int    `json:"game_rule" gorm:"game_rule"`               // -- 场次级别
	AmountTax      int64  `json:"amount_tax" gorm:"amount_tax"`             // 平台金币税收
	PlatformGameId int    `json:"platform_game_id" gorm:"platform_game_id"` // 游戏平台id
}

var r,_ = New("/Users/zonst/Desktop/workplace/awesomeProject/goTest/auto_project",AppCoinTaxReq{},AppCoinTax{},false)

func TestGenerateModel(t *testing.T) {
	t.Run("集成测试", func(t *testing.T) {
		// path 请求地址
		// reqSrc 请求结构体
		// src 数据库结构体
		// writeMyself 是否自己写sql
		err := GenerateData("/Users/zonst/Desktop/workplace/awesomeProject/goTest/auto_project",
			AppCoinTaxReq{}, AppCoinTax{},false)
		assert.Nil(t, err)
	})
}

func TestStruct(t *testing.T) {
	t.Run("struct format", func(t *testing.T) {
		formatPrint(r.Src)
	})
}



