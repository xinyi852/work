package bootstrap

import (
	"racent.com/pkg/config"
	"racent.com/pkg/database"
	"racent.com/pkg/logger"
	"time"
)

// SetupDB 初始化数据库和 ORM
func SetupDB() {
	opt := &database.Options{
		Username:           config.GetString("database.mysql.username"),
		Password:           config.GetString("database.mysql.password"),
		Host:               config.GetString("database.mysql.host"),
		Port:               config.GetInt("database.mysql.port"),
		DbName:             config.GetString("database.mysql.database"),
		Charset:            config.GetString("database.mysql.charset"),
		MaxOpenConnections: config.GetInt("database.mysql.max_open_connections"),
		MaxIdleConnections: config.GetInt("database.mysql.max_idle_connections"),
		MaxLifeSeconds:     time.Duration(config.GetInt("database.mysql.max_life_seconds")) * time.Second,
		Logger:             logger.Logger,
	}
	err := database.InitConnect(opt)
	if err != nil {
		logger.ErrorJSON(config.GetString("app.name"), "bootstrap", err.Error())
	} else {
		logger.InfoString(config.GetString("app.name"), "bootstrap", config.GetString("database.connection"))
	}
}
