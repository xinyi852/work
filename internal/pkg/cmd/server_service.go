package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"plesk/internal/pkg/bootstrap"
	"plesk/internal/pkg/global"
	"plesk/pkg/console"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"racent.com/pkg/config"
	"racent.com/pkg/logger"
)

// CmdAPIServer represents the available web sub-command.
var CmdAPIServer = &cobra.Command{
	Use:   "monitor-service-server",
	Short: "Start monitor service web server",
	Run:   runAPIService,
	Args:  cobra.NoArgs,
}

var apiServer *http.Server

func runAPIService(cmd *cobra.Command, args []string) {
	// 设置 gin 的运行模式，支持 debug, release, test
	// release 会屏蔽调试信息，官方建议生产环境中使用
	// 非 release 模式 gin 终端打印太多信息，干扰到我们程序中的 Log
	// 故此设置为 release，有特殊情况手动改为 debug 即可
	gin.SetMode(gin.ReleaseMode)

	// gin 实例
	router := gin.New()

	// 初始化路由绑定
	bootstrap.SetupRoute(router)

	apiServer = &http.Server{
		Addr:    ":" + config.GetString("app.port"),
		Handler: router,
	}

	go func() {
		logger.InfoString(config.GetString("app.name"), "Listening", config.GetString("app.port"))
		if err := apiServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.ErrorString("CMD", "api-server", err.Error())
			console.Exit("Unable to start api-server, error:" + err.Error())
		}
	}()

	signal.Notify(global.QuitChan, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	for {
		s := <-global.QuitChan
		switch s {
		case syscall.SIGHUP:
			logger.InfoString("CMD", "api-server", "收到终端断开信号，忽略")
		case os.Interrupt, syscall.SIGINT, syscall.SIGTERM:
			logger.InfoString("CMD", "api-server", "开始停止")
			shutdownAPIServer()
		}
	}
}

func shutdownAPIServer() {
	defer func() {
		logger.InfoString("CMD", "api-server", "Server exit")
		// 刷新日志
		err := logger.Sync()
		console.ExitIf(err)
		os.Exit(0)
	}()

	logger.InfoString("CMD", "api-server", "Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		logger.ErrorString("CMD", "api-server", "Server Shutdown:"+err.Error())
	}
}
