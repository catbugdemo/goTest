package auto_model

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"time"
)

// 支持 使用 go-cache 的二级缓存，

// 确定结构体
type struct_name struct {
	Id    int    `gorm:"id;default" json:"id" form:"id"`
	Str   string `gorm:"str;default" json:"str" form:"str"`
	Int   int    `gorm:"int;default" json:"int" form:"int"`
	Float float32
	Bool  bool
	Json  json.RawMessage
}

// 确定缓存名称 -- 单个
func (o struct_name) RedisKey() string {
	return fmt.Sprintf("struct_name:%d", o.Id)
}

// 确定缓存名称 -- 数组
func (o struct_name) ArrayRedisKey() string {
	return fmt.Sprintf("struct_name")
}

// 确定缓存时间范围
func (o struct_name) RedisSecondDuration() int {
	// TODO set redis 5-25 minute
	return int(time.Now().Unix()%25+5) * 60
}

// 确定数据库
func (o struct_name) DB() *gorm.DB {
	return nil
}

// 支持使用 redis 缓存 , 该缓存只 支持 string 类型的单个缓存
// 该种 string 类型每次都会有序列化和反小序列化的开销，
// 优点：简化编程，合理使用序列化可以提高内存使用效率
// 缺点：序列化和反序列化都有一定的开始

// GetFromRedis 获取缓存数据
func (o *struct_name) GetFromRedis(key string, conn redis.Conn) error {
	if key == "" {
		return errors.New("key is nil")
	}
	bytes, e := redis.Bytes(conn.Do("GET", key))
	if e == redis.ErrNil {
		return errors.New("not find in redis nor database")
	}
	if e != nil {
		return errors.WithStack(e)
	}
	if e = json.Unmarshal(bytes, &o); e != nil {
		return errors.WithStack(e)
	}
	return nil
}

// SyncToRedis 同步缓存 -1 为永久
func (o struct_name) SyncToRedis(key string, conn redis.Conn) error {
	if key == "" {
		return errors.New("key is nil")
	}
	buf, e := json.Marshal(o)
	if e != nil {
		return errors.WithStack(e)
	}
	if o.RedisSecondDuration() == -1 {
		if _, e = conn.Do("SET", key, buf); e != nil {
			return errors.WithStack(e)
		}
	} else {
		if _, e = conn.Do("SETEX", key, o.RedisSecondDuration(), buf); e != nil {
			return errors.WithStack(e)
		}
	}
	return nil
}

// DeleteFromRedis 删除缓存
func (o struct_name) DeleteFromRedis(key, arrayKey string, conn redis.Conn) error {
	if key != "" {
		if _, e := conn.Do("DEL", key); e != nil {
			return errors.WithStack(e)
		}
	}

	if arrayKey != "" {
		conn.Do("HDEL", arrayKey, o.Id)
	}

}
