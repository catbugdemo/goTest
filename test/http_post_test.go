package test

import (
	"bytes"
	"encoding/json"
	"github.com/fwhezfwhez/errorx"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var client = http.Client{

}

func HTTPPOSTForm(URL string, param map[string]string, rs interface{}) error {
	var values = url.Values{}
	for k, v := range param {
		values.Add(k, v)
	}

	rsp, e := client.PostForm(URL, values)

	if e != nil {
		return errorx.Wrap(e)
	}

	if rsp == nil {
		return errorx.NewServiceError("rsp nil", 5)
	}

	if rsp.Body == nil {
		return errorx.NewServiceError("rsp nil", 6)
	}

	if e := json.NewDecoder(rsp.Body).Decode(rs); e != nil {
		return errorx.Wrap(e)
	}
	return nil
}

type Param struct {
	Name string `json:"name"`
}


func TestPost(t *testing.T) {
	r := SetupRouter()

	w := httptest.NewRecorder()
	p := Param{
		Name: "123",
	}
	buf, _ = json.Marshal(p)

	req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(buf))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w,req)

	assert.Equal(t, 200,w.Code)
	m := map[string]string{
		"name": "123",
	}
	e := HTTPPOSTForm("http://localo/test", m, &p)
	assert.Nil(t, e)
}

// 封装获取结果

// 对 httptest 进行的 gin 封装测试
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/test", func(c *gin.Context) {
		var param Param
		if e := c.Bind(&param); e != nil {
			c.JSON(400, gin.H{"errCode": e})
			return
		}
		c.JSON(200, gin.H{"code": "200","msg":param.Name})
	})
	return r
}

