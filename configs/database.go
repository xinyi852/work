package configs

import "racent.com/pkg/config"

func init() {
	config.Add("database", func() map[string]interface{} {
		return map[string]interface{}{
			// 是否启用数据库
			"enable_db": config.Env("db.enable", false),

			// 默认数据库
			"connection": config.Env("db.connection", "mysql"),

			"mysql": map[string]interface{}{

				// 数据库连接信息
				"host":     config.Env("db.host", "127.0.0.1"),
				"port":     config.Env("db.port", "3306"),
				"database": config.Env("db.database", ""),
				"username": config.Env("db.username", ""),
				"password": config.Env("db.password", ""),
				"charset":  "utf8mb4",

				// 连接池配置
				"max_idle_connections": config.Env("db.max_idle_connections", 25),
				"max_open_connections": config.Env("db.max_open_connections", 100),
				"max_life_seconds":     config.Env("db.max_life_seconds", 5*60),
			},

			"sqlite": map[string]interface{}{
				"database": config.Env("db.sql_file", "database/database.db"),
			},
		}
	})
}
