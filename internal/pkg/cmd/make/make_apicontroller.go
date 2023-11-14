package make

import (
	"fmt"
	"plesk/internal/pkg/file"
	"plesk/pkg/console"
	"strings"

	"github.com/spf13/cobra"
)

var CmdMakeAPIController = &cobra.Command{
	Use:   "apicontroller",
	Short: "Create api controller，exmaple: make apicontroller v1/auth/login",
	Run:   runMakeAPIController,
	Args:  cobra.ExactArgs(1), // 只允许且必须传 1 个参数
}

func runMakeAPIController(cmd *cobra.Command, args []string) {

	// 处理参数，要求附带 API 版本（v1 或者 v2）
	array := strings.Split(args[0], "/")
	arrayLen := len(array)
	if arrayLen < 2 {
		console.Exit("api controller name format: v1/user")
	}

	// apiVersion 用来拼接目标路径
	// name 用来生成 cmd.Model 实例
	apiVersion, name := array[0], array[arrayLen-1]
	model := makeModelFromString(name)

	// 文件目录
	dirPath := "internal/api/controllers/"
	if arrayLen > 1 {
		for i := 0; i < arrayLen-1; i++ {
			dirPath += array[i] + "/"
		}
	}

	// 如果目录不存在，就创建
	err := file.MkdirAll(dirPath, 0750)
	console.ExitIf(err)

	// 组建目标目录
	fileName := fmt.Sprintf("%s_controller.go", model.PackageName)

	// 基于模板创建文件（做好变量替换）
	createFileFromStub(dirPath+fileName, "apicontroller", model, map[string]string{"{{version}}": apiVersion})
}
