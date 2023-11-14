package bootstrap

import (
	"racent.com/pkg/config"
	"racent.com/pkg/logger"
)

// SetupLogger 初始化 Logger
func SetupLogger() {
	opt := &logger.Options{
		FileAddr:  config.GetString("log.filename"),
		MaxSize:   config.GetInt("log.max_size"),
		MaxBackup: config.GetInt("log.max_backup"),
		MaxAge:    config.GetInt("log.max_age"),
		Compress:  config.GetBool("log.compress"),
		LogType:   config.GetString("log.type"),
		Level:     config.GetString("log.level"),
		Env:       config.GetString("app.env"),
	}
	logger.InitLogger(opt)
}
