package middlewares

import (
	"bytes"
	"io"
	"plesk/internal/pkg/apisign"
	"plesk/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"racent.com/pkg/config"
)

func APIVerifySign() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.GetBool("app.debug") {
			err := verifyAPISignature(c)
			if err != nil {
				response.JsonFail(c, response.ErrApiSignatureInvalid, err.Error())
				return
			}
		}

		c.Next()
	}
}

func verifyAPISignature(c *gin.Context) error {
	// 获取GET请求数据
	requestQuery := make(map[string]string)
	for k, _ := range c.Request.URL.Query() {
		requestQuery[k] = c.Query(k)
	}

	// 获取body请求数据
	var requestBody []byte
	if c.Request.Body != nil {
		requestBody, _ = io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
	}

	// 验证签名
	apiClient := apisign.Client{Gin: c}
	return apiClient.VerifySignature(c.Request.Method, requestQuery, string(requestBody))
}
