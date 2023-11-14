package configs

import "racent.com/pkg/config"

func init() {
	config.Add("cert_driver", func() map[string]interface{} {
		return map[string]interface{}{
			"list": config.Env("cert_driver", nil),
		}
	})
}
