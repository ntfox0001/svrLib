package util

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
)

// 这个函数只能用来调试，不可用于正式功能
func GetGID() uint64 {
	fmt.Println("this function is only for debug!!")
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
