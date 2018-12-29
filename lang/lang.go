package lang

import (
	"testapi/tools"
)

// 传入语言类型和key
func GetLang(lang string) *JsonData {
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

	// 返回JSON数据
	return &res
}
