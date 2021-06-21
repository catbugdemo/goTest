package auto_generate

import (
	"fmt"
	"testing"
)

func TestGenerateController(t *testing.T) {
	//前端是否有参数传入
	hasParam := true
	rs := GenerateController(hasParam, map[string]string{
		//自定义controller模块名
		"${handler_name}": "ReceiveAward",
		//service层连接名称
		"${handler_service}": "timeAwardService.ReceiveAward",
		//自定义异常
		"${handle_error}": "timeAwardUtil.SaveError(errorx.Wrap(e))",
	})
	_ = rs
	fmt.Println(rs)
}
