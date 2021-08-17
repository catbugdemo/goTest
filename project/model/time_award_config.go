package model

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"main.go/project/pkg/setting"
	"main.go/project/utils"
	"time"
)

type TimeAwardConfig struct {
	Id        int       `gorm:"column:id;default:" json:"id" form:"id"`
	CreatedAt time.Time `gorm:"column:created_at;default:" json:"created_at" form:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;default:" json:"updated_at" form:"updated_at"`

	GameId     int    `gorm:"column:game_id;default:" json:"game_id" form:"game_id"`
	OrderIndex int    `gorm:"column:order_index;default:" json:"order_index" form:"order_index"`
	AppChannel string `gorm:"column:app_channel;default:'vx'" json:"app_channel" form:"app_channel"`

	Props     json.RawMessage `gorm:"column:props;default:" json:"props" form:"props"`
	StartHour int             `gorm:"column:start_hour;default:" json:"start_hour" form:"start_hour"`
	EndHour   int             `gorm:"column:end_hour;default:" json:"end_hour" form:"end_hour"`
}

func (o TimeAwardConfig) DB() *gorm.DB {
	return setting.Db
}

// 设置缓存
var TimeAwardConfigRedisKeyFormat = "time_award_config:%d:%d:%d"

func (o TimeAwardConfig) RedisKey() string {
	return fmt.Sprintf(TimeAwardConfigRedisKeyFormat,o.GameId,o.OrderIndex,o.AppChannel)
}

// 设置数组缓存(如果需要用到)
var ArrayTimeAwardConfigRedisKeyFormat = "time_award_config:%d:%d"

func (o TimeAwardConfig) ArrayRedisKey() string {
	return fmt.Sprintf(TimeAwardConfigRedisKeyFormat,o.GameId,o.OrderIndex)
}

// RedisSecondDuration 缓存时间
func (o TimeAwardConfig) RedisSecondDuration() int {
	// TODO set its redis duration, default 1-7 day,  return -1 means no time limit
	return int(time.Now().Unix()%7+1) * 60*60*24
}

// MustGet 获取信息，当前只有缓存和数据库
// 获取过程
// 1.先从缓存中获取
// 2.再从数据库中获取
func (o *TimeAwardConfig) MustGet(conn redis.Conn,engine *gorm.DB) error {
	//从缓存中获取
	buf, e := utils.GetFromRedis(o.RedisKey(), conn)
	if e!= nil && e.Error() == "not found redis nor db"{
		return e
	}

	//未从缓存中获取
	if e!= nil{
		var count int
		 if e = engine.Count(&count).Error;e!=nil{
		 	return e
		 }
		 if count == 0 {
		 	if o.RedisSecondDuration() == -1 {
				if _, e = conn.Do("SET", o.RedisKey(), "NX");e!=nil{
					return e
				}
			}else {
				if _, e = conn.Do("SET", o.RedisKey(),"EX",o.RedisSecondDuration(), "NX");e!=nil{
					return e
				}
			}
		 }

		 //查询数据库
		if e = engine.First(&o).Error;e!=nil{
			return e
		}

		//同步缓存
		if e = utils.SyncToRedis(o.RedisKey(), o.RedisSecondDuration(), o, conn);e!=nil{
			return e
		}
		return nil
	}

	//将缓存绑定--(太远了)
	if e = json.Unmarshal(buf, &o);e!=nil{
		return e
	}

	return nil
}


