package utils

import (
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
	"time"
)

func RedisSecondDuration() int {
	// TODO set its redis duration, default 1-7 day,  return -1 means no time limit
	return int(time.Now().Unix()%7+0) * 60 * 60 * 24
}

// SyncToRedis 同步缓存
// 1.判断key是否为为空
// 2.判断 time 是否为空 (time.Unix)
//   (1)空(-1)设置永久缓存
//   (2)设置时间缓存(时间必须为一段时间内的随机数)
// 3.同步缓存
// 4.判断操作是否有错误
func SyncToRedis(redisKey string, time int, value interface{}, conn redis.Conn) error {
	if redisKey == "" {
		return errors.New("redis key is nil")
	}

	buf, e := json.Marshal(value)
	if e != nil {
		return e
	}

	if time == -1 {
		if _, e := conn.Do("SET", redisKey, buf);e!=nil{
			return e
		}
	} else {
		if _,e := conn.Do("SETEX",redisKey,time,buf);e!=nil{
			return e
		}
	}

	return nil
}

// GetFromRedis 获取缓存
// 推荐使用 -- func(o *Test) GetFromRedis () error -- 减少 json.Unmarshal(buf,&o) 操作
// 1.判断key是否为为空
// 2.获取缓存( byte 格式)
//	(1)判断操作是否有误
//  (2)判断缓存值value是否为空
//  (3)判断缓存值是否为 DISABLE (防止缓存穿透，默认将value 设置为 DISABLE)
func GetFromRedis(redisKey string,conn redis.Conn) ([]byte,error) {
	if redisKey == ""{
		return nil, errors.New("redis key is nil")
	}

	buf, e := redis.Bytes(conn.Do("GET", redisKey))

	if e == redis.ErrNil {
		return nil, e
	}

	if e != nil && e!= redis.ErrNil {
		return nil,e
	}

	if string(buf) == "DISABLE" {
		return nil, errors.New("not found in db nor redis")
	}

	return buf,nil
}

// DeleteFromRedis 删除缓存 (当有数组缓存时，单个缓存更新数据缓存必须删除重新获取)
// 1.判断缓存key是否为空
// 2.删除缓存
// 3.判断数组缓存key是否为空
// 4.删除缓存
// 可同时作用于ArrayDeleteFromRedis
func DeleteFromRedis(redisKey,arryRedisKey string,conn redis.Conn) error {
	if redisKey == ""{
		if _, e := conn.Do("DEL", redisKey);e!=nil{
			return e
		}
	}

	if arryRedisKey == ""{
		if _, e := conn.Do("DEL", arryRedisKey);e!=nil{
			return e
		}
	}

	return nil
}
