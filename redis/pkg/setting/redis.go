package setting

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

var (
	RedisPool *redis.Pool
	conn redis.Conn
)


func InitRedisQuickLink() {
	var err error
	conn, err = redis.Dial("tcp", Setting.Redis.Host)
	if err != nil {
		conn.Close()
		fmt.Println("fail to dial redis:", err)
	}
}

func InitRedisPool()  {
	RedisPool = &redis.Pool{
		//最大闲置连接
		MaxIdle: Setting.Redis.MaxIdle,
		//最大活动数
		MaxActive: Setting.Redis.MaxActive,
		//数据库连接
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", Setting.Redis.Host)
			if err != nil {
				c.Close()
				fmt.Printf("fail to dial redis: %v\n", err)
				return nil, err
			}
			//密码认证
			if Setting.Redis.Password != "" {
				if _, err = c.Do("AUTH", Setting.Redis.Password); err != nil {
					c.Close()
					fmt.Printf("fail to auth redis: %v\n", err)
					return nil, err
				}
			}
			//redis 缓存数据库认证
			if _, err = c.Do("SELECT", Setting.Redis.Db); err != nil {
				c.Close()
				fmt.Printf("fail to SELECT DB redis: %v\n", err)
				return nil, err
			}
			return c, err
		},
		//测试连接是否正常
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				c.Close()
				fmt.Printf("fail to ping redis: %v\n", err)
				return err
			}
			return nil
		},
	}
}



func GetRedisQuick() *redis.Conn {
	return &conn
}
