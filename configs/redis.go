package configs

import "racent.com/pkg/config"

func init() {
	config.Add("redis", func() map[string]interface{} {
		return map[string]interface{}{
			// 是否启用redis
			"enable_redis": config.Env("redis.enable", false),
			"host":         config.Env("redis.host", "127.0.0.1"),
			"port":         config.Env("redis.port", "6379"),
			"password":     config.Env("redis.password", ""),

			// 业务类存储使用 1 (图片验证码、短信验证码、会话)
			"database": config.Env("redis.main_db", 1),
		}
	})
}
