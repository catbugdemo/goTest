package setting

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

var (
	RedisPool *redis.Pool
	conn redis.Conn
)

type redisConfig struct {
	Redis
}

func InitRedisQuickLink() {
	r := redisConfig{
		Redis: Setting.Redis,
	}

	var err error
	conn, err = redis.Dial("tcp", r.Host)
	if err != nil {
		conn.Close()
		fmt.Println("fail to dial redis:", err)
	}
}

func InitRedisPool()  {
	r := redisConfig{
		Redis: Setting.Redis,
	}

	RedisPool = &redis.Pool{
		//最大闲置连接
		MaxIdle: r.MaxIdle,
		//最大活动数
		MaxActive: r.MaxActive,
		//数据库连接
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", r.Host)
			if err != nil {
				c.Close()
				fmt.Printf("fail to dial redis: %v\n", err)
				return nil, err
			}
			//密码认证
			if r.Password != "" {
				if _, err = c.Do("AUTH", r.Password); err != nil {
					c.Close()
					fmt.Printf("fail to auth redis: %v\n", err)
					return nil, err
				}
			}
			//redis 缓存数据库认证
			if _, err = c.Do("SELECT", r.Db); err != nil {
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
	log.Println("connect redis success")

}



func GetRedisQuick() *redis.Conn {
	return &conn
}
