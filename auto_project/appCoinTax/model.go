package appCoinTax

//根据请求结构体生成
type AppCoinTaxModuleReq struct {
	StartTime    string `json:"start_date"`   // 日期
	EndTime      string `json:"end_date"`     // 用户平台id
	GameId       int    `json:"game_id"`      // --子游戏id
	GameAreaId   int    `json:"game_area_id"` // --登陆渠道
	AppChannel   string `json:"app_channel"`  // -- 场次级别
	LoginChannel string `json:"login_channel"`
	IsDownload   bool   `json:"is_download"`
}

// 数据库表结构体
type AppCoinTax struct {
	CalDate        string `json:"cal_date" gorm:"cal_date"`           // 日期
	GameId         int    `json:"game_id" gorm:"game_id"`             // 用户平台id
	GameAreaId     int    `json:"game_area_id" gorm:"game_area_id"`   // --子游戏id
	AppChannel     string `json:"app_channel" gorm:"app_channel"`     // --登陆渠道
	LoginChannel   string `json:"login_channel" gorm:"login_channel"` // -- 场次级别
	GameRule       int    `json:"game_rule" gorm:"game_rule"`
	AmountTax      int64  `json:"amount_tax" gorm:"amount_tax"`
	PlatformGameId int    `json:"platform_game_id" gorm:"platform_game_id"`
}

// 开发人员自己写sql
func GetAppCoinTaxQueryList(model AppCoinTaxModuleReq) ([]AppCoinTax, int, error) {
	a := make([]AppCoinTax, 0)
	var count int

	if model.LoginChannel == "all" {

		sql := "select * from table_name where id=? and name=?"
		if err := db.DataDB.Table("app_coin_tax").Raw(sql, (model.GameAreaId * model.GameId)).Find(&a).Error; err != nil {
			return a, count, err
		}

		if err := db.DataDB.Table("app_coin_tax").
			Where("cal_date between ? and ?", model.StartTime, model.EndTime).
			Where("game_id=?", model.GameId).
			Where("game_area_id=?", model.GameAreaId).
			Where("app_channel=?", model.AppChannel).Group("call_date,game_id,game_area_id,app_channel").Count(&count).Error; err != nil {
			return a, count, err
		}
		return a, count, nil

	} else {

		sql := "select * from table_name where id=? and name=?"
		if err := db.DataDB.Table("app_coin_tax").Raw(sql, model.GameAreaId).Find(&a).Error; err != nil {
			return a, count, err
		}

		if err := db.DataDB.Table("app_coin_tax").
			Where("cal_date between ? and ?", model.StartTime, model.EndTime).
			Where("game_id=?", model.GameId).
			Where("game_area_id=?", model.GameAreaId).
			Where("app_channel=?", model.AppChannel).
			Where("login_channel=?", model.LoginChannel).Count(&count).Error; err != nil {
			return a, count, err
		}
		return a, count, nil

	}

}
