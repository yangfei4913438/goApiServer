package conf

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"testapi/tools"
)

func init() {
	setLog()
}

func setLog() {
	//日志配置
	var cfg = map[string]interface{}{
		//保存的文件名
		"filename": tools.GetRootPath() + beego.AppConfig.String("log_filename"),
		//level 日志保存的时候的级别，默认是 Trace 级别, 数字越小打印的范围越小，日志级别越高
		//0-emergency, 1-[emergency、alert]，2-[emergency、alert、critical，error也不打印],
		//3-[不打印warning、notice、debug以及info日志], 生产环境打印级别设置为3
		//4-[不打印notice、debug以及info日志], 5-[不打印debug和info日志], 6-[不打印debug日志], 7-[全部级别]
		"level": beego.AppConfig.DefaultInt("log_level", 7),
		//是否按照每天log rotate，默认是 true
		"daily": beego.AppConfig.DefaultBool("log_daily", true),
		//文件最多保存多少天，默认保存 7 天
		"maxdays": beego.AppConfig.DefaultInt("log_maxdays", 7),
		//每个文件保存的最大行数，默认值 1000000
		"maxlines": beego.AppConfig.DefaultInt("log_maxlines", 0), //0也是默认值的一种写法
		//每个文件保存的最大尺寸，默认值是 1 << 28, //256 MB
		"maxsize": beego.AppConfig.DefaultInt("log_maxsize", 0), //0也是默认值的一种写法
	}
	// 转化日志为字符串
	log_conf, _ := json.Marshal(cfg)
	// 支持命令行显示日志
	beego.SetLogger(logs.AdapterConsole, "console")
	// 支持日志打印到文件中
	beego.SetLogger(logs.AdapterFile, string(log_conf))
}

/*
日志级别，从低到高。示例代码：
beego.Debug("beego test debug log")
beego.Info("beego test info log")
beego.Notice("beego test notice log")
beego.Warning("beego test warning log")
beego.Error("beego test error log")
beego.Critical("beego test critical log")
beego.Alert("beego test alert log")
beego.Emergency("beego test emergency log")
*/
