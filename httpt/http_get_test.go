package httpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"testing"
)

type clientMsg struct {
	GameId     int    `json:"game_id"`
	UserId     int    `json:"user_id"`
	AppChannel string `json:"app_channel"`
	AppId      string `json:"app_id"`
	Version    int    `json:"version"`
}


func TestPostGetToken(t *testing.T) {
	t.Run("集成测试:", func(t *testing.T) {
		token, e := PostGetToken("http://xyx.zonst.com/dev/api/auth/generate-token/", "application/json", client)

		assert.Nil(t, e)
		assert.NotNil(t, token)
		fmt.Println(token)
	})
}

func TestGetToken(t *testing.T) {
	t.Run("集成测试:GetToken", func(t *testing.T) {
		//测试使用方法
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
}


func TestShowData(t *testing.T) {
	t.Run("集成测试：不传入数据", func(t *testing.T) {
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

	t.Run("集成测试：传入数据", func(t *testing.T) {
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
}

func TestCheckData(t *testing.T) {
	t.Run("集成测试：CheckData", func(t *testing.T) {
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

func TestStringToJson(t *testing.T) {
	t.Run("测试：stringToJson", func(t *testing.T) {
		str := "{\"data\":{\"list\":[{\"time_award_config_id\":2,\"grey\":true,\"order_index\":1,\"props\":[{\"prop_id\":19,\"prop_num\":10000,\"expire_in\":-1}],\"start_time\":\"7:00\",\"end_time\":\"9:00\",\"vstate\":1},{\"time_award_config_id\":3,\"grey\":true,\"order_index\":2,\"props\":[{\"prop_id\":19,\"prop_num\":20000,\"expire_in\":-1}],\"start_time\":\"10:00\",\"end_time\":\"11:00\",\"vstate\":1},{\"time_award_config_id\":4,\"grey\":true,\"order_index\":3,\"props\":[{\"prop_id\":19,\"prop_num\":30000,\"expire_in\":-1}],\"start_time\":\"16:00\",\"end_time\":\"17:00\",\"vstate\":1}]},\"tip\":\"success\",\"tip_id\":0}"
		body, e := json.Marshal(str)
		assert.Nil(t, e)

		var prettyJSON bytes.Buffer
		e = json.Indent(&prettyJSON, body, "", "\t")
		if e != nil {
			log.Println("fail to indent:", e)
		}

		fmt.Println(string(prettyJSON.Bytes()))
	})

}
