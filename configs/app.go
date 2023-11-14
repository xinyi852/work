package configs

import "racent.com/pkg/config"

func init() {
	config.Add("app", func() map[string]interface{} {
		return map[string]interface{}{
			// 应用名称
			"name": config.Env("app.name", "plesk_service"),

			// 当前环境，用以区分多环境，一般为 local, stage, production, test
			"env": config.Env("app.env", "production"),

			// 是否进入调试模式
			"debug": config.Env("app.debug", false),

			// 应用服务端口
			"port": config.Env("app.port", "8056"),

			// 加密会话、JWT 加密
			"key": config.Env("app.key", "33cfa0ea06016a6532b96da32f304a446a9df0de"),

			// 用以生成链接
			"url": config.Env("app.url", "http://localhost:3000"),

			// 设置时区，JWT 里会使用，日志记录里也会使用到
			"timezone": config.Env("app.timezone", "Asia/Shanghai"),

			// 语言包
			"locale": config.Env("app.locale", "en"),
			// 回调次数间隔时间分钟
			"webhook_max_times": config.Env("app.webhook_max_times"),
		}
	})
}
