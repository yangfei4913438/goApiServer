package lang

import (
	"testapi/tools"
)

// 下面的函数被执行后，就可以通过这个全局变量，在全局获取语言数据
var CurrLang *JsonData

// 传入语言类型
func GetLang(lang string) {
	// 定义接收数据的对象变量
	var res JsonData

	// 定义json文件的路径变量
	var filePath string

	// 指定JSON语言文件
	switch lang {
	case "zh-cn":
		// 匹配简体中文
		filePath = tools.GetRootPath() + "lang/zh-cn.json"
	case "zh-tw":
		// 匹配繁体中文
		filePath = tools.GetRootPath() + "lang/zh-tw.json"
	default:
		// 默认匹配英文
		filePath = tools.GetRootPath() + "lang/en.json"
	}

	// 解析JSON语言文件
	tools.ParseJsonFile(filePath, &res)

	// 赋值给全局变量
	CurrLang = &res
}
