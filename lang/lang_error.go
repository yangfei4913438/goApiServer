package lang

import (
	"errors"
	"github.com/astaxie/beego"
	"goApiServer/tools"
	"strings"
)

// 错误结构体，前端如果使用多语言，则需要自己准备多语言模板，如果不使用多语言，则可以直接显示后端的错误消息。
type errorMessage struct {
	Code    int
	Message string
}

// 消息可以转换成错误类型
func (err *errorMessage) ErrorType() error {
	return errors.New(err.Message)
}

type RRSError struct {
	ERR500  *errorMessage
	Err1001 *errorMessage
	Err1002 *errorMessage
	Err1003 *errorMessage
	Err1004 *errorMessage
	Err1005 *errorMessage
	Err1006 *errorMessage
	Err1007 *errorMessage
	Err1008 *errorMessage
}

var RRSErrorInfo RRSError

var cnLang jsonData
var twLang jsonData
var enLang jsonData

// 获取环境变量
var env = beego.AppConfig.String("runmode")

func init() {
	if env == "prod" {
		absPath, _ := tools.GetRootPath()
		var path = *absPath
		// 解析JSON语言文件, 这里统一加载后，每次请求接口，就不会再去加载了。
		// 所有的语言文件，都在这里进行加载。
		_ = tools.ParseJsonFile(path+"lang/zh-cn.json", &cnLang)
		_ = tools.ParseJsonFile(path+"lang/zh-tw.json", &twLang)
		_ = tools.ParseJsonFile(path+"lang/en.json", &enLang)
	} else {
		// 解析JSON语言文件, 这里统一加载后，每次请求接口，就不会再去加载了。
		// 所有的语言文件，都在这里进行加载。
		_ = tools.ParseJsonFile("lang/zh-cn.json", &cnLang)
		_ = tools.ParseJsonFile("lang/zh-tw.json", &twLang)
		_ = tools.ParseJsonFile("lang/en.json", &enLang)
	}

	// 默认设置为英文
	SetLang("en")
}

func SetLang(lang string) {
	switch strings.ToLower(lang) {
	case "zh-cn":
		CurrLang = &cnLang
	case "zh-tw":
		CurrLang = &twLang
	default:
		CurrLang = &enLang
	}

	// 系统错误

	RRSErrorInfo.ERR500 = &errorMessage{Code: 500, Message: CurrLang.ErrInfo.Err500}

	// 用户错误

	RRSErrorInfo.Err1001 = &errorMessage{Code: 1001, Message: CurrLang.ErrInfo.Err1001}
	RRSErrorInfo.Err1002 = &errorMessage{Code: 1002, Message: CurrLang.ErrInfo.Err1002}
	RRSErrorInfo.Err1003 = &errorMessage{Code: 1003, Message: CurrLang.ErrInfo.Err1003}
	RRSErrorInfo.Err1004 = &errorMessage{Code: 1004, Message: CurrLang.ErrInfo.Err1004}
	RRSErrorInfo.Err1005 = &errorMessage{Code: 1005, Message: CurrLang.ErrInfo.Err1005}
	RRSErrorInfo.Err1006 = &errorMessage{Code: 1006, Message: CurrLang.ErrInfo.Err1006}
	RRSErrorInfo.Err1007 = &errorMessage{Code: 1007, Message: CurrLang.ErrInfo.Err1007}
	RRSErrorInfo.Err1008 = &errorMessage{Code: 1008, Message: CurrLang.ErrInfo.Err1008}
}
