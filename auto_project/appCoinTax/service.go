package appCoinTax

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "time"
)



func AppCoinTaxQueryList(c *gin.Context) {

    defer logstat.LogStat("AppCoinTaxQueryList", c, time.Now())
  	// 请求结构体
	type AppCoinTaxReq struct {
    StartTime  string    `json:"start_date"` // 日期  
    EndTime  string    `json:"end_date"` // 用户平台id  
    GameId  int    `json:"game_id"` // --子游戏id  
    GameAreaId  int    `json:"game_area_id"` // --登陆渠道  
    AppChannel  string    `json:"app_channel"` // -- 场次级别  
    LoginChannel  string    `json:"login_channel"` 
    IsDownload  bool    `json:"is_download"` 
 
}
  
    req := AppCoinTaxReq{}
  
    if err := c.ShouldBindJSON(&req); err != nil {
        logging.Errorf("AppCoinTaxReq,err:%v\n", err.Error())
        c.JSON(http.StatusOK, gin.H{"error": -1, "message": err.Error()})
        return
    }
  
  // model 层函数 get + 当前方法名
  
    data, count, err := GetAppCoinTaxQueryList(AppCoinTaxModuleReq{
		StartTime:    req.StartTime,
		EndTime:    req.EndTime,
		GameId:    req.GameId,
		GameAreaId:    req.GameAreaId,
		AppChannel:    req.AppChannel,
		LoginChannel:    req.LoginChannel,
		IsDownload:    req.IsDownload,

	})
  
    if err != nil {
        logging.Errorf("GetAppCoinTaxList,err:%v\n", err.Error())
        c.JSON(http.StatusOK, gin.H{"error": -1, "message": err.Error()})
        return
    }
  	
	
if !req.IsDownload {
	c.JSON(http.StatusOK, gin.H{"error": 0, "message": "", "data": data, "count": count})
	return
} else {
	TableTitle := []string{"日期", "用户平台id", "--子游戏id", "--登陆渠道", "-- 场次级别"} 		// 数据库结构注释
	var TableDate [][]string	// 数据

	for _, v := range data {
		TableDate = append(TableDate, []string{
			v.CalDate,
			strconv.Itoa(int(v.GameId)),
			strconv.Itoa(int(v.GameAreaId)),
			v.AppChannel,
			v.LoginChannel,
			strconv.Itoa(int(v.GameRule)),
			strconv.FormatInt(v.AmountTax,10),
			strconv.Itoa(int(v.PlatformGameId)),

		}) //  根据数据库表结构 和数据的类型进行数据区分 转换成 string
	}

	//调用方法  表名参数
	utils.ExportToExcel(c, TableTitle, TableDate, "123")

}

}
