package cmd

import (
	"encoding/pem"
	"fmt"
	"github.com/spf13/cobra"
	"racent.com/pkg/logger"
	"regexp"
)

var CmdPlay = &cobra.Command{
	Use:   "play",
	Short: "专门调试的方法，调试完成后，记住清除测试代码",
	Run:   runPlay,
}

// 专门调试的方法，调试完成后，记住清除测试代码
func runPlay(cmd *cobra.Command, args []string) {

	certChain := `-----BEGIN CERTIFICATE-----
AAAAAA
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
BBBBBB
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
CCCCCC
-----END CERTIFICATE-----`

	re := regexp.MustCompile(`(?s)-----BEGIN CERTIFICATE-----(.*?)-----END CERTIFICATE-----`)

	matches := re.FindAllString(certChain, -1)

	var certBlocks []*pem.Block

	for _, match := range matches {
		fmt.Println(match)
		block, _ := pem.Decode([]byte(match))
		fmt.Println(block)
		if block != nil {
			certBlocks = append(certBlocks, block)
		}
	}
	logger.Dump(certBlocks)
	//for _, match := range matches {
	//	fmt.Println(match)
	//	fmt.Println("")
	//	fmt.Println("")
	//	fmt.Println("")
	//}
}
