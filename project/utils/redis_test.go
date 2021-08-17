package utils

import (
	"encoding/json"
	"github.com/agiledragon/gomonkey"
	"github.com/alicebob/miniredis/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func getMiniRedis() redis.Conn {
	s, e := miniredis.Run()
	if e != nil {
		log.Fatalln("Fail to run miniredis:",e)
	}

	conn, e := redis.Dial(`tcp`, s.Addr())
	if e != nil {
		log.Fatalln("Fail to dial redis:",e)
	}

	return conn
}


func TestSyncToRedis(t *testing.T) {
	t.Run("单元测试：SyncToRedis", func(t *testing.T) {
		jsonFunc := gomonkey.ApplyFunc(json.Marshal, func(v interface{}) ([]byte, error) {
			return nil, nil
		})
		defer jsonFunc.Reset()

		conn := getMiniRedis()
		defer conn.Close()


		e1 := SyncToRedis("123", 123,"123", conn)
		e2 := SyncToRedis("-1", -1,"-1", conn)

		assert.Nil(t, e1)
		assert.Nil(t, e2)
	})

	t.Run("测试：miniredis", func(t *testing.T) {
		s, e := miniredis.Run()
		if e != nil {
			log.Fatalln("Fail to run miniredis:",e)
		}

		c, e := redis.Dial(`tcp`, s.Addr())
		defer c.Close()
		if e != nil {
			log.Fatalln("Fail to dial redis:",e)
		}

		assert.NotNil(t, c)

	})

}

func TestGetFromRedis(t *testing.T) {
	t.Run("单元测试:GetFromRedis: not found in db nor redis", func(t *testing.T) {
		c := getMiniRedis()
		defer c.Close()

		c.Do("SET","1","DISABLE")

		_, e := GetFromRedis("1", c)

		assert.Equal(t, "not found in db nor redis",e.Error())
	})

	t.Run("单元测试:GetFromRedis:success", func(t *testing.T) {
		c := getMiniRedis()
		defer c.Close()
		c.Do("SETEX","1",123,"1")
		fromRedis, e := GetFromRedis("1", c)

		assert.Nil(t, e)
		assert.NotNil(t, fromRedis)
	})

}

func TestDeleteFromRedis(t *testing.T) {
	t.Run("单元/连条测试:DeleteFromRedis", func(t *testing.T) {
		conn := getMiniRedis()
		defer conn.Close()

		conn.Do("SET","1","1")
		conn.Do("SET","2","2")

		e := DeleteFromRedis("1", "2", conn)

		assert.Nil(t, e)
	})
}