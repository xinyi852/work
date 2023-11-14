// Package file 文件操作辅助函数
package file

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Put 将数据存入文件
func Put(data []byte, to string) error {
	err := os.WriteFile(to, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Exists 判断文件是否存在
func Exists(fileToCheck string) bool {
	if _, err := os.Stat(fileToCheck); os.IsNotExist(err) {
		return false
	}
	return true
}

// 创建多级目录，类似 mkdir -p
func MkdirAll(dirPath string, perm os.FileMode) error {
	err := os.MkdirAll(dirPath, perm)
	if err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

func FileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func CreateFile(fileName string) (*os.File, error) {
	if Exists(fileName) {
		rmErr := os.RemoveAll(fileName)
		if rmErr != nil {
			return nil, rmErr
		}
	}
	f, createErr := os.Create(fileName)

	if createErr != nil {
		return nil, createErr
	}
	return f, nil
}

// 读取文件内容
func ReadFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var (
		// 注意buf不能够声明为没有长度的切片,因为官方文档说了最终返回的字节数目n可能小于len(buf), 所以切片需要长度否则读取不了数据
		buf  = make([]byte, 4096)
		data []byte
	)
	reader := bufio.NewReader(f)
	for {
		n, readErr := reader.Read(buf)
		if readErr != nil && readErr != io.EOF {
			return nil, readErr
		}
		// 文件末尾那么退出
		if n == 0 {
			break
		}
		// 将读取到的数据追加到data切片中
		data = append(data, buf[:n]...)
	}
	return data, nil
}
