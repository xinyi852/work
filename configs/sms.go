package configs

import "racent.com/pkg/config"

func init() {
	config.Add("sms", func() map[string]interface{} {
		return map[string]interface{}{

			// 默认是阿里云的测试 sign_name 和 template_code
			"aliyun": map[string]interface{}{
				"access_key_id":     config.Env("sms.aliyun.access_id"),
				"access_key_secret": config.Env("sms.aliyun.access_secret"),
				"sign_name":         config.Env("sms.aliyun.sign_name", "阿里云短信测试"),
				"template_code":     config.Env("sms.aliyun.template_code", "SMS_154950909"),
			},

			// 腾讯云
			"tencent": map[string]interface{}{
				"access_key_id":     config.Env("sms.tencent.access_id"),
				"access_key_secret": config.Env("sms.tencent.access_secret"),
				"sign_name":         config.Env("sms.tencent.sign_name"),
				"template_code":     config.Env("sms.tencent.template_code"),
				"sdk_app_id":        config.Env("sms.tencent.sdk_app_id"),
			},
			// 百度云

		}
	})
}
