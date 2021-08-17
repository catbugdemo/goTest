package auto_generate

import (
	"fmt"
	"testing"
)

func TestGenerateController(t *testing.T) {

	rs := GenerateController(map[string]interface{}{
		//自定义controller模块名
		"${handler_name}": "ReceiveAward",
		//service层连接名称
		"${handler_service}": "timeAwardService.ReceiveAward",
		//自定义异常
		"${handle_error}": "timeAwardUtil.SaveError(errorx.Wrap(e))",
		// 是否有用户权限认证
		"has_authMiddleware": false,
		//前端是否有参数传入
		"has_param": false,
	})
	_ = rs
	fmt.Println(rs)
}
