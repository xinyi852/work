package bootstrap

import (
	"racent.com/pkg/config"
	"racent.com/pkg/i18n"
)

func SetupLocale() {
	lang := config.GetString("app.locale")
	i18n.InitConfig(lang, []string{"etc/lang/en.toml", "etc/lang/zh-CN.toml"})
}
