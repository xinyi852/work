package configs

import "racent.com/pkg/config"

func init() {
	config.Add("user_default_quota", func() map[string]interface{} {
		return map[string]interface{}{
			"list": config.Env("user_default_quota", nil),
		}
	})
}
