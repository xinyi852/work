// Package bootstrap 处理程序初始化逻辑
package bootstrap

import (
	"net/http"
	"plesk/internal/api/middlewares"
	"plesk/internal/routes"
	"strings"

	"github.com/gin-gonic/gin"
	"racent.com/pkg/config"
	"racent.com/pkg/logger"
)

// SetupRoute 路由初始化
func SetupRoute(router *gin.Engine) {
	// 注册全局中间件
	registerGlobalMiddleWare(router)

	logger.InfoString(config.GetString("app.name"), "bootstrap", "route")
	// 注册 API 路由
	routes.RegisterAPIRoutes(router)

	//  配置 404 路由
	setup404Handler(router)
}

func registerGlobalMiddleWare(router *gin.Engine) {
	router.Use(
		middlewares.Logger(),
		middlewares.Recovery(),
		middlewares.Cors(),
	)
}

func setup404Handler(router *gin.Engine) {
	// 处理 404 请求
	router.NoRoute(func(c *gin.Context) {
		// 获取标头信息的 Accept 信息
		acceptString := c.Request.Header.Get("Accept")
		if strings.Contains(acceptString, "text/html") {
			// 如果是 HTML 的话
			c.String(http.StatusNotFound, "页面返回 404")
		} else {
			// 默认返回 JSON
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"data":    nil,
				"message": "路由未定义，请确认 url 和请求方法是否正确。",
				"errors":  nil,
			})
		}
	})
}
