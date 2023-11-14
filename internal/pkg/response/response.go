// Package response 响应处理工具
package response

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"racent.com/pkg/logger"
)

type Result struct {
	Code    int                 `json:"code"`
	Data    interface{}         `json:"data"`
	Message string              `json:"message"`
	Errors  map[string][]string `json:"errors"`
}

// JsonSuccess 业务成功 响应 200 和 对应数据
func JsonSuccess(c *gin.Context, data interface{}) {
	ret := Result{
		Code:    Success,
		Data:    data,
		Message: CodeText(Success),
	}
	c.JSON(http.StatusOK, ret)
}

// JsonFail 业务失败 响应 200 和 对应状态信息
func JsonFail(c *gin.Context, responseCode int, msg ...string) {

	ret := Result{
		Code:    responseCode,
		Message: defaultMessage(CodeText(responseCode), msg...),
	}
	c.AbortWithStatusJSON(http.StatusOK, ret)
}

// 没有权限时候使用
// Abort401 响应 401，未传参 msg 时使用默认消息
func Abort401(c *gin.Context, msg ...string) {
	ret := Result{
		Code:    ErrUnauthorized,
		Message: defaultMessage(CodeText(ErrUnauthorized), msg...),
	}
	c.AbortWithStatusJSON(http.StatusOK, ret)
}

// Abort404 响应 200，未传参 msg 时使用默认消息
func Abort404(c *gin.Context, msg ...string) {
	ret := Result{
		Code:    ErrNotFound,
		Message: defaultMessage(CodeText(ErrNotFound), msg...),
	}
	c.AbortWithStatusJSON(http.StatusOK, ret)
}

// Abort422 响应 422，未传参 msg 时使用默认消息
func Abort422(c *gin.Context, msg ...string) {
	ret := Result{
		Code:    ErrUnprocessableEntity,
		Message: defaultMessage(CodeText(ErrUnprocessableEntity), msg...),
	}
	c.AbortWithStatusJSON(http.StatusOK, ret)
}

// 全部接口响应 200
// 处理请求时出现错误 err，会附带返回 error 信息，如登录错误、找不到 ID 对应的 Model
func Error(c *gin.Context, err error, msg ...string) {
	logger.LogIf(err)

	// error 类型为『数据库未找到内容』
	if err == gorm.ErrRecordNotFound {
		Abort404(c)
		return
	}

	ret := Result{
		Code:    ErrDefault,
		Message: defaultMessage(err.Error(), msg...),
	}

	c.AbortWithStatusJSON(http.StatusOK, ret)
}

// ValidationError 处理表单验证不通过的错误，返回的 JsonSuccess 示例：
//
//	{
//	    "errors": {
//	        "phone": [
//	            "手机号为必填项，参数名称 phone",
//	            "手机号长度必须为 11 位的数字"
//	        ]
//	    },
//	    "message": "请求验证不通过，具体请查看 errors"
//	}
func ValidationError(c *gin.Context, errors map[string][]string) {
	ret := Result{
		Code:    ErrValidate,
		Message: CodeText(ErrValidate),
		Errors:  errors,
	}
	c.AbortWithStatusJSON(http.StatusOK, ret)
}

// defaultMessage 内用的辅助函数，用以支持默认参数默认值
// Go 不支持参数默认值，只能使用多变参数来实现类似效果
func defaultMessage(defaultMsg string, msg ...string) (message string) {
	if len(msg) > 0 {
		message = msg[0]
	} else {
		message = defaultMsg
	}
	return message
}

func Format(code int) *responseCode {
	return &responseCode{code, CodeText(code)}
}
