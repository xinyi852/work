// Package routes 注册路由
package routes

import (
	v1 "plesk/internal/api/controllers/v1"

	// "plesk/internal/api/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterAPIRoutes(r *gin.Engine) {
	// 监控服务 v1Group 版本的路由都将存放到这里
	// v1Group := r.Group("/v1", middlewares.APIVerifySign())
	v1Group := r.Group("/v1")
	{

		pleskCtrl := new(v1.HostingController)
		pleskServerCtrl := new(v1.PleskServerController)
		pleskServerGroupCtrl := new(v1.PleskServerGroupController)
		pleskProductCtrl := new(v1.PleskProductController)
		pleskGroup := v1Group.Group("hosting")
		{
			pleskGroup.POST("/server_group/create", pleskServerGroupCtrl.Create)
			pleskGroup.POST("/server_group/delete", pleskServerGroupCtrl.Delete)
			pleskGroup.POST("/server_group/update", pleskServerGroupCtrl.Update)
			pleskGroup.POST("/server_group/list", pleskServerGroupCtrl.List)

			pleskGroup.POST("/plan/create", pleskProductCtrl.Create)
			pleskGroup.POST("/plan/delete", pleskProductCtrl.Delete)
			pleskGroup.POST("/plan/update", pleskProductCtrl.Update)
			pleskGroup.POST("/plan/list", pleskProductCtrl.List)

			pleskGroup.POST("/service/create", pleskCtrl.Create)
			pleskGroup.POST("/service/upgrade", pleskCtrl.Upgrade)
			pleskGroup.POST("/service/change_password", pleskCtrl.ChangePassword)
			pleskGroup.POST("/service/sign_login", pleskCtrl.SignLogin)
			pleskGroup.POST("/service/suspend", pleskCtrl.Suspend)
			pleskGroup.POST("/service/active", pleskCtrl.Active)
			pleskGroup.POST("/service/delete", pleskCtrl.Delete)
			pleskGroup.POST("/service/info", pleskCtrl.Info)

			pleskGroup.POST("/server/create", pleskServerCtrl.Create)
			pleskGroup.POST("/server/list", pleskServerCtrl.ServerList)
			pleskGroup.POST("/server/delete", pleskServerCtrl.Delete)
			pleskGroup.POST("/server/update", pleskServerCtrl.Update)

		}
	}
}
