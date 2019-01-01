package dbs

import (
	"github.com/astaxie/beego"
	redis "github.com/yangfei4913438/redis-full"
	"strings"
	"time"
)

//redis对外接口
var RedisDB redis.RedisCache

func initRedis() {
	hosts := beego.AppConfig.String("redis.host")
	password := beego.AppConfig.DefaultString("redis.password", "")
	database := beego.AppConfig.DefaultInt("redis.db", 0)
	MaxIdle := beego.AppConfig.DefaultInt("redis.maxidle", 100)
	MaxActive := beego.AppConfig.DefaultInt("redis.maxactive", 1000)
	IdleTimeout := time.Duration(beego.AppConfig.DefaultInt("redis.idletimeout", 600)) * time.Second

	//通过赋值对外接口来使用
	RedisDB = redis.NewRedisCache(hosts, password, database, MaxIdle, MaxActive, IdleTimeout, 24*time.Hour)

	if err := RedisDB.CheckRedis(); err != nil {
		panic("Redis Server:" + hosts + " Connect failed: " + err.Error() + "!")
	} else {
		beego.Info("Connect Redis Server(" + hosts + ") to successful!")
	}
}

// 删掉包含指定前缀的key
func CleanRedis(prefixes ...string) error {
	// 获取所有的key
	keys, err := RedisDB.Keys()
	if err != nil {
		beego.Error(err)
		return err
	}

	// 遍历这些key，删掉包含指定前缀的key
	for _, key := range keys {
		for _, prefix := range prefixes {
			if strings.HasPrefix(key, prefix) {
				if err2 := RedisDB.Del(key); err2 != nil {
					beego.Error(err2)
					return err2
				}
				// 删除完成了，就不用继续检查了，换下个KEY
				break
			}
		}
	}

	return nil
}
