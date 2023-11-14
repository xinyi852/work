package global

import "os"

// QuitChan 退出通道
var QuitChan = make(chan os.Signal, 1)
