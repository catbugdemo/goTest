package auto_generate

import (
	"strings"
)

//GenerateController
//传入数据:被替代的名称
// ${handler_name} --- 模块名称
// ${handler_service} --- service层模块名
// ${handler_error} --- 自定义异常抛出类型类型
func GenerateController(replacement map[string]interface{}) string {
	handleDefaultController(replacement)

	var fro_control_mould = `
// Auto Generate controller

//${handler_name}
func ${handler_name}(c *gin.Context) {
`
	var has_authMiddleware = `
	//获取用户数据
	pl, e := authMiddleware.GetClientPayload(c)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"tip_id": 1, "tip": "token参数异常，请稍后重试", "debug_message": errorx.Wrap(e).Error()})
		return
	}`
	var cen_control_mould = `
	//根据前端是否传入数据来确定是否删除
	type Param struct {
		// TODO set its  front end data
	}
	var param Param
	if e = c.Bind(&param); e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"tip_id": 1, "tip": "success", "debug_message": errorx.Wrap(e).Error()})
		return
	}
	// ---------
`

	var aft_control_mould = `
	//service层操作
	//配置缓存
	conn := redistool.RedisPool.Get()
	defer conn.Close()
	// TODO modify its parameter
	data, e := ${handler_service}(pl.GameId, pl.UserId, , conn)
	if e != nil {
		if se, ok := errorx.IsServiceErr(e); ok {
			c.JSON(http.StatusOK, gin.H{"tip_id": se.Errcode, "tip": se.Errmsg})
			return
		}
		${handle_error}
		c.JSON(http.StatusOK, gin.H{"tip_id": 2,"tip":"${handler_service} 出现异常", "debug_message": errorx.Wrap(e).Error()})
		return
	}
	
	c.JSON(http.StatusOK,gin.H{"tip_id":0,"tip":"success","data":data})
}
`
	var result = fro_control_mould
	if replacement["has_authMiddleware"].(bool) {
		result += has_authMiddleware
	}

	if replacement["has_param"].(bool) {
		result += cen_control_mould + aft_control_mould
	} else {
		result += aft_control_mould
	}

	result = strings.ReplaceAll(result, "${handler_name}", replacement["${handler_name}"].(string))
	result = strings.ReplaceAll(result, "${handler_service}", replacement["${handler_service}"].(string))
	result = strings.ReplaceAll(result, "${handle_error}", replacement["${handle_error}"].(string))

	return result
}

func handleDefaultController(replacement map[string]interface{}) {
	if replacement["${handler_name}"] == "" {
		replacement["${handler_name}"] = "DefaultController"
	}

	if replacement["${handler_service}"] == "" {
		replacement["${handler_service}"] = "defaultService.DefaultService"
	}

	if replacement["${handle_error}"] == "" {
		replacement["${handle_error}"] = "defaultUtil.SaveError(errorx.Wrap(e))"
	}

}
