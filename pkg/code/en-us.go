package code

var enUSText = map[int]string{
	ServerError:           "Internal server error",
	TooManyRequests:       "Too many requests",
	ParamBindError:        "Parameter error",
	AuthorizationError:    "Authorization error",
	UrlSignError:          "URL signature error",
	CacheSetError:         "Failed to set cache",
	CacheGetError:         "Failed to get cache",
	CacheDelError:         "Failed to del cache",
	CacheNotExist:         "Cache does not exist",
	ResubmitError:         "Please do not submit repeatedly",
	HashIdsEncodeError:    "HashID encryption failed",
	HashIdsDecodeError:    "HashID decryption failed",
	RBACError:             "No access",
	RedisConnectError:     "Failed to connection Redis",
	MySQLConnectError:     "Failed to connection MySQL",
	WriteConfigError:      "Failed to write configuration file",
	SendEmailError:        "Failed to send mail",
	MySQLExecError:        "SQL execution failed",
	GoVersionError:        "Go Version mismatch",
	SocketConnectError:    "Socket not connected",
	SocketSendError:       "Socket message sending failed",
	ParamDnsTypeBindError: "Parameter dns_type error",
	ParamMxBindError:      "Parameter mx error,require:1-50",
	ParamTTLBindError:     "Parameter ttl error,require:600-60400",
	ParamURLHTTPBindError: "Parameter record value error,need prefix: http:// https://",

	AuthorizedCreateError:    "Failed to create caller",
	AuthorizedListError:      "Failed to get caller list",
	AuthorizedDeleteError:    "Failed to delete caller",
	AuthorizedUpdateError:    "Failed to update caller",
	AuthorizedDetailError:    "Failed to get caller details",
	AuthorizedCreateAPIError: "Failed to create caller API address",
	AuthorizedListAPIError:   "Failed to get caller API address list",
	AuthorizedDeleteAPIError: "Failed to delete caller API address",

	AdminCreateError:             "Failed to create administrator",
	AdminListError:               "Failed to get administrator list",
	AdminDeleteError:             "Failed to delete administrator",
	AdminUpdateError:             "Failed to update administrator",
	AdminResetPasswordError:      "Reset password failed",
	AdminLoginError:              "Login failed",
	AdminLogOutError:             "Exit failed",
	AdminModifyPasswordError:     "Failed to modify password",
	AdminModifyPersonalInfoError: "Failed to modify personal information",
	AdminMenuListError:           "Failed to get administrator menu authorization list",
	AdminMenuCreateError:         "Administrator menu authorization failed",
	AdminOfflineError:            "Offline administrator failed",
	AdminDetailError:             "Failed to get personal information",

	MenuCreateError:       "Failed to create menu",
	MenuUpdateError:       "Failed to update menu",
	MenuDeleteError:       "Failed to delete menu",
	MenuListError:         "Failed to get menu list",
	MenuDetailError:       "Failed to get menu details",
	MenuCreateActionError: "Failed to create menu action",
	MenuListActionError:   "Failed to get menu action list",
	MenuDeleteActionError: "Failed to delete menu action",

	CronCreateError:  "Failed to create cron",
	CronUpdateError:  "Failed to update menu",
	CronListError:    "Failed to get cron list",
	CronDetailError:  "Failed to get cron detail",
	CronExecuteError: "Failed to execute cron",

	DomainCreateErr:   "Domain name creation failed. The domain name has been created",
	DomainListErr:     "Domain not exist,please check",
	RecordListErr:     "Record not exist,please check",
	RecordModifyErr:   "Record modify err",
	DomainIsCreateErr: "Domain name creation failed. The domain name has been created by yourself",
	RecordIsExistErr:  "The record already exists, please do not add it again",
	RecordDeleteErr:   "The record not exist,please check",
	RecordCreateErr:  "The record create fail,please check value",
}
