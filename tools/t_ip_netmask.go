package tools

import (
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
)

// 16进制子网掩码，其中两位字符，转成10进制字符
func str16To10(str string) string {
	n, err := strconv.ParseUint(str, 16, 32)
	if err != nil {
		beego.Error("子网掩码转换出错: " + err.Error())
	}
	return fmt.Sprint(n)
}

// 传入一个 16 进制子网掩码字符串，转换成 10 进制的子网掩码字符串
func NetMask16to10(str string) string {
	var s1 = str16To10(str[:2])
	var s2 = str16To10(str[2:4])
	var s3 = str16To10(str[4:6])
	var s4 = str16To10(str[6:])
	return s1 + "." + s2 + "." + s3 + "." + s4
}
