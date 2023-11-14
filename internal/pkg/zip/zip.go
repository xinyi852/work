package zip

import (
	"archive/zip"
	"bytes"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Zip 压缩文件或目录
// @params dst io.Writer 压缩文件可写流
// @params src string    待压缩源文件/目录路径
func Zip(dst io.Writer, src string) error {
	// 强转一下路径
	src = filepath.Clean(src)
	// 提取最后一个文件或目录的名称
	baseFile := filepath.Base(src)
	// 判断src是否存在
	_, err := os.Stat(src)
	if err != nil {
		return err
	}

	// 通文件流句柄创建一个ZIP压缩包
	zw := zip.NewWriter(dst)
	// 延迟关闭这个压缩包
	defer zw.Close()

	// 通过filepath封装的Walk来递归处理源路径到压缩文件中
	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		// 是否存在异常
		if err != nil {
			return err
		}

		// 通过原始文件头信息，创建zip文件头信息
		zfh, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// 赋值默认的压缩方法，否则不压缩
		zfh.Method = zip.Deflate

		// 移除绝对路径
		tmpPath := path
		index := strings.Index(tmpPath, baseFile)
		if index > -1 {
			tmpPath = tmpPath[index:]
		}
		// 替换文件名，并且去除前后 "\" 或 "/"
		tmpPath = strings.Trim(tmpPath, string(filepath.Separator))
		// 替换一下分隔符，zip不支持 "\\"
		zfh.Name = strings.ReplaceAll(tmpPath, "\\", "/")
		// 目录需要拼上一个 "/" ，否则会出现一个和目录一样的文件在压缩包中
		if info.IsDir() {
			zfh.Name += "/"
		}

		// 写入文件头信息，并返回一个ZIP文件写入句柄
		zfw, err := zw.CreateHeader(zfh)
		if err != nil {
			return err
		}

		// 仅在他是标准文件时进行文件内容写入
		if zfh.Mode().IsRegular() {
			// 打开要压缩的文件
			sfr, err := os.Open(path)
			if err != nil {
				return err
			}
			defer sfr.Close()

			// 将srcFileReader拷贝到zipFilWrite中
			_, err = io.Copy(zfw, sfr)
			if err != nil {
				return err
			}
		}

		// 搞定
		return nil
	})
}

// Unzip 解压压缩文件
// @params dst string      解压后的目标路径
// @params src *zip.Reader 压缩文件可读流
func Unzip(dst string, src *zip.Reader) error {
	// 强制转换一遍目录
	dst = filepath.Clean(dst)

	// 遍历压缩文件
	for _, file := range src.File {
		// 在闭包中完成以下操作可以及时释放文件句柄
		err := func() error {
			// 跳过文件夹
			if file.Mode().IsDir() {
				return nil
			}

			// 配置输出目标路径
			filename := filepath.Join(dst, file.Name)
			// 创建目标路径所在文件夹
			e := os.MkdirAll(filepath.Dir(filename), os.ModeDir)
			if e != nil {
				return e
			}

			// 打开这个压缩文件
			zfr, e := file.Open()
			if e != nil {
				return e
			}
			defer zfr.Close()

			// 创建目标文件
			fw, e := os.Create(filename)
			if e != nil {
				return e
			}
			defer fw.Close()

			// 执行拷贝
			_, e = io.Copy(fw, zfr)
			if e != nil {
				return e
			}

			// 拷贝成功
			return nil
		}()

		// 是否发生异常
		if err != nil {
			return err
		}
	}

	// 解压完成
	return nil
}

// ZipToFile 压缩至文件
// @params dst string 压缩文件目标路径
// @params src string 待压缩源文件/目录路径
// @return     error  错误信息
func ZipToFile(dst, src string) error {
	// 创建一个ZIP文件
	fw, err := os.Create(filepath.Clean(dst))
	if err != nil {
		return err
	}
	defer fw.Close()

	// 执行压缩
	return Zip(fw, src)
}

// ZipToBytes 压缩至字节流
// @params src string 待压缩源文件/目录路径
// @return []byte 压缩字节流
// @return error  错误信息
func ZipToBytes(src string) ([]byte, error) {
	// 创建一个缓冲区
	buf := bytes.NewBuffer(nil)

	// 执行压缩
	err := Zip(buf, src)
	if err != nil {
		return nil, err
	}

	// 返回缓冲区的二进制流
	return buf.Bytes(), nil
}

// UnzipFromFile 解压压缩文件
// @params dst string 解压后目标路径
// @params src string 压缩文件目标路径
func UnzipFromFile(dst, src string) error {
	// 打开压缩文件
	zr, err := zip.OpenReader(filepath.Clean(src))
	if err != nil {
		return err
	}
	defer zr.Close()

	// 解压
	return Unzip(dst, &zr.Reader)
}

// UnzipFromBytes 解压压缩字节流
// @params dst string 解压后目标路径
// @params src []byte 压缩字节流
func UnzipFromBytes(dst string, src []byte) error {
	// 通过字节流创建zip的Reader对象
	zr, err := zip.NewReader(bytes.NewReader(src), int64(len(src)))
	if err != nil {
		return err
	}

	// 解压
	return Unzip(dst, zr)
}
