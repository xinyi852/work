package configs

import "racent.com/pkg/config"

func init() {
	config.Add("monitor_certificate", func() map[string]interface{} {
		return map[string]interface{}{
			// 是否启用监控，默认不启用
			"enable": config.Env("monitor_certificate.enable", false),
			// 时间间隔 单位分钟
			"interval_time": config.Env("monitor_certificate.interval_time", 1),
		}
	})
}
