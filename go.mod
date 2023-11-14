module plesk

go 1.20

require (
	github.com/gertd/go-pluralize v0.2.1
	github.com/gin-gonic/gin v1.8.1
	github.com/google/uuid v1.3.0
	github.com/iancoleman/strcase v0.2.0
	github.com/jordan-wright/email v4.0.1-0.20210109023952-943e75fe5223+incompatible
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d
	github.com/spf13/cast v1.5.0
	github.com/spf13/cobra v1.6.0
	github.com/thedevsaddam/govalidator v1.9.10
	go.uber.org/zap v1.24.0
	golang.org/x/crypto v0.3.0
	gorm.io/gorm v1.25.1
	racent.com/pkg/cert v0.0.1
	racent.com/pkg/config v0.0.1
	racent.com/pkg/database v0.0.1
	racent.com/pkg/helpers v0.0.1
	racent.com/pkg/i18n v0.0.1
	racent.com/pkg/logger v0.0.1
	racent.com/pkg/redis v0.0.1
)

replace (
	racent.com/pkg/cert => ../pkg/cert
	racent.com/pkg/config => ../pkg/config
	racent.com/pkg/database => ../pkg/database
	racent.com/pkg/helpers => ../pkg/helpers
	racent.com/pkg/httpclient => ../pkg/httpclient
	racent.com/pkg/i18n => ../pkg/i18n
	racent.com/pkg/logger => ../pkg/logger
	racent.com/pkg/redis => ../pkg/redis
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.11.0 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/goccy/go-json v0.9.11 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/nicksnyder/go-i18n/v2 v2.2.1 // indirect
	github.com/pelletier/go-toml/v2 v2.0.6 // indirect
	github.com/redis/go-redis/v9 v9.0.4 // indirect
	github.com/spf13/afero v1.9.3 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.15.0 // indirect
	github.com/subosito/gotenv v1.4.2 // indirect
	github.com/ugorji/go/codec v1.2.7 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/net v0.4.0 // indirect
	golang.org/x/sys v0.3.0 // indirect
	golang.org/x/text v0.5.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/mysql v1.5.1 // indirect
	racent.com/pkg/httpclient v0.0.1 // indirect

)
