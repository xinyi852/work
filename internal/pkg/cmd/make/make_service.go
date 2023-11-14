package make

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var CmdMakeService = &cobra.Command{
	Use:   "service",
	Short: "Crate service file, example: make service user",
	Run:   runMakeService,
	Args:  cobra.ExactArgs(1), // 只允许且必须传 1 个参数
}

func runMakeService(cmd *cobra.Command, args []string) {

	// 格式化模型名称，返回一个 Model 对象
	model := makeModelFromString(args[0])

	// 确保模型的目录存在，例如 `internal/models/`
	dir := fmt.Sprintf("internal/services/%s/", model.PackageName)
	// os.MkdirAll 会确保父目录和子目录都会创建，第二个参数是目录权限，使用 0777
	os.MkdirAll(dir, os.ModePerm)

	// 替换变量
	createFileFromStub(dir+model.PackageName+"_service.go", "service", model)
}
