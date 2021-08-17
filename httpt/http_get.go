package httpt

import (
	"bytes"
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"log"
	"net/http"
)

type TokenGet struct {
	Data  string `json:"data" binding:"required"`
	Exp   int    `json:"exp" binding:"required"`
	TipId int    `json:"tip_id" binding:"required"`
}

// PostGetToken 以POST方式获取token
// 1.json解析
// 2.发送post请求
// 3.读取io流
// 4.转为结构体
func PostGetToken(url string,contentType string,value interface{}) (TokenGet,error) {
	buf, e := json.Marshal(&value)
	if e != nil {
		log.Println("fail to json Marshal")
		return TokenGet{}, e
	}

	post, e := http.Post(url, contentType, bytes.NewBuffer(buf))
	if e != nil {
		log.Fatalln("fail to post:", url)
		return TokenGet{}, e
	}
	//读取io流
	all, e := ioutil.ReadAll(post.Body)
	if e != nil {
		log.Fatalln("fail to read Body")
		return TokenGet{}, e
	}
	//转为结构体
	get := TokenGet{}
	e = json.Unmarshal(all, &get)
	if e != nil {
		log.Fatalln("fail to json unmarshal")
		return TokenGet{}, e
	}

	return get, nil

}

// GetToken 获取token
// 1.发送请求
// 2.读取io流
// 3.绑定格式
// 使用方法：
//	buf, e := json.Marshal(client)
//	assert.Nil(t, e)
//	req, e := http.NewRequest("POST", "http://xyx.zonst.com/dev/api/auth/generate-token/", bytes.NewReader(buf))
//	req.Header.Set("Content-Type", "application/json")
//	assert.Nil(t, e)
// 	token, e := GetToken(req)
func GetToken(req *http.Request) (string, error) {
	client := &http.Client{}

	resp, e := client.Do(req)
	if e != nil {
		log.Println("fail to client do")
		return "",e
	}
	defer resp.Body.Close()

	//读取io流
	buf, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		log.Println("fail to read body")
		return "", e
	}

	//绑定格式
	get := TokenGet{}
	e = json.Unmarshal(buf, &get)
	if e != nil {
		log.Fatalln("fail to json unmarshal")
		return "", e
	}

	return get.Data, nil
}

// ShowData 展示数据
// 1.转换格式
// 2.自定义配置请求
// 3.发送请求
// 4.转换数据
func ShowData(method ,url string,header map[string]string,value interface{}) (string,error) {
	body, e := sendRequest(method, url, header, value)
	if e != nil {
		return "", e
	}

	var prettyJSON bytes.Buffer
	e = json.Indent(&prettyJSON, body, "", "\t")
	if e != nil {
		log.Println("fail to indent:",e)
		return "", e
	}

	return string(prettyJSON.Bytes()),nil
}


// CheckData 确定数据(封装了sendRequest+simpleJson)
func CheckData(method ,url string,header map[string]string,value interface{}) (*simplejson.Json,error) {
	body, e := sendRequest(method, url, header, value)
	if e != nil {
		return nil, e
	}

	js, e := simplejson.NewJson(body)
	return js,nil
}


func sendRequest(method ,url string,header map[string]string,value interface{}) ([]byte,error)   {
	buf, e := json.Marshal(&value)
	if e != nil {
		log.Println("fail to json marshal:",e)
		return nil, e
	}

	req, e := http.NewRequest(method, url, bytes.NewReader(buf))
	if e != nil {
		log.Println("fail to http NewRequest:",e)
		return nil, e
	}
	for k,v :=range header{
		req.Header.Add(k,v)
	}

	client := &http.Client{}
	resp, e := client.Do(req)
	if e != nil {
		log.Println("fail to client do:",e)
		return nil, e
	}
	defer resp.Body.Close()

	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		log.Println("fail to read:",e)
		return nil, e
	}
	return body, nil
}