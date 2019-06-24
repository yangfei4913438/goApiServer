package tools

import "unsafe"

// 字节数组 转 字符串
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
