package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var json_str string = `{"rc" : 0,
  "error" : "Success",
  "type" : "stats",
  "progress" : 100,
  "job_status" : "COMPLETED",
  "result" : {
    "total_hits" : 803254,
    "starttime" : 1528434707000,
    "endtime" : 1528434767000,
    "fields" : [ ],
    "timeline" : {
      "interval" : 1000,
      "start_ts" : 1528434707000,
      "end_ts" : 1528434767000,
      "rows" : [ {
        "start_ts" : 1528434707000,
        "end_ts" : 1528434708000,
        "number" : "x12887"
      }, {
        "start_ts" : 1528434720000,
        "end_ts" : 1528434721000,
        "number" : "x13028"
      }, {
        "start_ts" : 1528434721000,
        "end_ts" : 1528434722000,
        "number" : "x12975"
      }, {
        "start_ts" : 1528434722000,
        "end_ts" : 1528434723000,
        "number" : "x12879"
      }, {
        "start_ts" : 1528434723000,
        "end_ts" : 1528434724000,
        "number" : "x13989"
      } ],
      "total" : 803254
    },
      "total" : 8
  }
}`
// StringToJson string以json格式化输出
func TestStringToJson(t *testing.T) {
	//因为json格式化只允许[]bytes，所以所有的类型都要转为[]butes
	t.Run("stringtoJson", func(t *testing.T) {
		var prettyJson bytes.Buffer
		e := json.Indent(&prettyJson, []byte(json_str), " ", "\t")
		assert.Nil(t, e)
		fmt.Println(string(prettyJson.Bytes()))

	})
}