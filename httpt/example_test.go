package httpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var client = clientMsg{
	GameId:     66,
	UserId:     446823,
	AppChannel: "vx",
	AppId:      "asda",
	Version:    1,
}

func TestExample(t *testing.T) {
	t.Run("PostGetToken:", func(t *testing.T) {
		token, e := PostGetToken("http://xyx.zonst.com/dev/api/auth/generate-token/", "application/json", client)

		assert.Nil(t, e)
		assert.NotNil(t, token)
		fmt.Println(token)
	})

	t.Run("GetToken", func(t *testing.T) {
		buf, e := json.Marshal(client)
		assert.Nil(t, e)

		req, e := http.NewRequest("POST", "http://xyx.zonst.com/dev/api/auth/generate-token/", bytes.NewReader(buf))
		req.Header.Set("Content-Type", "application/json")
		assert.Nil(t, e)

		token, e := GetToken(req)
		assert.Nil(t, e)
		assert.NotNil(t, token)
		fmt.Println(token)
	})

	//可以作为项目的一般测试 -- 不传入数据
	t.Run("ShowData", func(t *testing.T) {
		//获取token
		token, e := PostGetToken("http://xyx.zonst.com/dev/api/auth/generate-token/", "application/json", client)
		assert.Nil(t, e)

		header := make(map[string]string)
		header["client-auth"] = token.Data

		//显示数据
		data, e := ShowData("POST", "http://xyx.zonst.com/dev/api/time-award-GetTimeAwardProcAndConfig/", header, nil)
		assert.Nil(t, e)

		fmt.Println(data)
	})
	//传入数据
	t.Run("ShowData", func(t *testing.T) {
		token, e := PostGetToken("http://xyx.zonst.com/dev/api/auth/generate-token/", "application/json", client)
		assert.Nil(t, e)

		header := make(map[string]string)
		header["client-auth"] = token.Data
		header["Content-Type"] = "application/json"

		type Param struct {
			ConfigId int `json:"config_id" form:"config_id" binding:"exists"`
		}

		param := Param{
			ConfigId: 2,
		}

		data, e := ShowData("POST", "http://xyx.zonst.com/dev/api/time-award-ReceiveAward/", header, param)
		assert.Nil(t, e)

		fmt.Println(data)
	})

	t.Run("CheckData", func(t *testing.T) {
		token, e := PostGetToken("http://xyx.zonst.com/dev/api/auth/generate-token/", "application/json", client)
		assert.Nil(t, e)

		header := make(map[string]string)
		header["client-auth"] = token.Data

		//确认数据
		js, e := CheckData("POST", "http://xyx.zonst.com/dev/api/time-award-GetTimeAwardProcAndConfig/", header, nil)
		assert.Nil(t, e)

		b, e := js.Get("data").Get("list").GetIndex(0).Get("grey").Bool()
		assert.Nil(t, e)
		assert.Equal(t, true,b)
	})
}





