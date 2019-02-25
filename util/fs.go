package util

import "os"

// 检查文件是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	// if os.IsNotExist(err) {
	// 	return false
	// }
	return false
}

// 创建目录
func CreatePath(path string) error {
	if !PathExists(path) {
		err := os.MkdirAll(path, os.ModePerm)
		return err
	}
	return nil
}
