package bootstrap

import (
	"fmt"
	"racent.com/pkg/config"
	"racent.com/pkg/logger"
	"racent.com/pkg/redis"
)

// SetupRedis 初始化 Redis
func SetupRedis() {
	// 建立 redis 连接
	opt := &redis.Options{
		Addr:     fmt.Sprintf("%v:%v", config.GetString("redis.host"), config.GetString("redis.port")),
		Username: config.GetString("redis.username"),
		Password: config.GetString("redis.password"),
		Db:       config.GetInt("redis.database"),
	}

	err := redis.ConnectRedis(opt)
	if err != nil {
		logger.ErrorString(config.GetString("app.name"), "redis Error", err.Error())
	} else {
		logger.InfoString(config.GetString("app.name"), "bootstrap", "redis")
	}

}
