package main

import (
	"fmt"
	"os"
	"plesk/configs"
	"plesk/internal/pkg/bootstrap"
	"plesk/internal/pkg/cmd"
	cmdMakePkg "plesk/internal/pkg/cmd/make"
	"plesk/pkg/console"

	"github.com/spf13/cobra"
	"racent.com/pkg/config"
)

func main() {
	//应用的主入口，默认调用 cmd.CmdSensorServer 命令
	var rootCmd = &cobra.Command{
		Use:   "monitor",
		Short: "monitor service",
		Long:  `Default will run "monitor" command, you can use "-h" flag to see all subcommands`,
		// rootCmd 的所有子命令都会执行以下代码
		PersistentPreRun: func(command *cobra.Command, args []string) {
			// 配置初始化，依赖命令行 --env 参数
			config.InitConfig(cmd.Env, "yaml")
			configs.Initialize()
			// 初始化 Logger
			bootstrap.SetupLogger()
			// 初始化语言包
			bootstrap.SetupLocale()
			// 初始化验证规则
			bootstrap.SetupValidators()
			// 根据配置，决定是否初始化数据库
			enableDB := config.GetBool("database.enable_db")
			if enableDB {
				// 初始化数据库
				bootstrap.SetupDB()
			}
			// 初始化 Redis
			// bootstrap.SetupRedis()

			// 初始化证书驱动
			// bootstrap.SetupCertDriver()

			// 初始化监控证书状态
			// bootstrap.SetupMonitorCertificate()

		},
	}

	// 注册子命令
	rootCmd.AddCommand(
		cmd.CmdAPIServer,
		cmd.CmdKey,
		cmdMakePkg.CmdMake,
		cmd.CmdMigrate,
		cmd.CmdPlay,
	)

	// 配置默认运行 Web 服务
	cmd.RegisterDefaultCmd(rootCmd, cmd.CmdAPIServer)

	// 注册全局参数，--env
	cmd.RegisterGlobalFlags(rootCmd)
	// 执行主命令
	if err := rootCmd.Execute(); err != nil {
		console.Exit(fmt.Sprintf("Failed to run app with %v: %s", os.Args, err.Error()))
	}
}
