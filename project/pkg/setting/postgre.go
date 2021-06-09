package setting

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"time"
)

var Db *gorm.DB

type gormConfig struct {
	Postgre
}

func InitDatabase() {
	g := gormConfig{
		Postgre: Setting.Postgre,
	}

	dbConfig := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=%s password=%s",
		g.Host,
		g.User,
		g.Name,
		"disable",
		g.Password,
	)

	db, e := gorm.Open(g.Type, dbConfig)
	if e != nil {
		log.Fatal("Fail to connect 'database':", g.Type, ". err:", e)
	}

	//设置表前缀
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return g.TablePrefix + defaultTableName
	}

	//关闭复数形式
	db.SingularTable(true)
	db.LogMode(true)
	Db = db
	log.Println("connect postgre success")

	// 自动重连，每60秒ping一次，失败时自动重连，重连间隔依次为3s,3s,15s,30s,60s,60s,60s.....
	// 未进行测试
	go func(dbConfig string) {
		var intervals = []time.Duration{3 * time.Second, 3 * time.Second, 15 * time.Second, 30 * time.Second, 60 * time.Second}
		for {
			time.Sleep(60 * time.Second)
			// Ping失败需要重新连接
			if e = Db.DB().Ping();e!=nil{
				for i := 0; i < len(intervals); i++ {
					db, e = gorm.Open(g.Type, dbConfig)
					if e != nil {
						fmt.Println("errors is:",e)
						time.Sleep(intervals[i])
						if i == len(intervals)- 1{
							i--
							continue
						}
					}

					gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
						return g.TablePrefix + defaultTableName
					}
					//关闭复数形式
					db.SingularTable(true)
					db.LogMode(true)
					Db = db
				}
			}
		}
	}(dbConfig)
}
