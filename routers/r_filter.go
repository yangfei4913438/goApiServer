package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"goApiServer/controllers"
	"goApiServer/dbs"
	"goApiServer/lang"
)

func RouterFilter() {
	// 过滤器函数，检查操作
	var check = func(ctx *context.Context) {
		// 语言处理
		lg := ctx.Request.Header.Get("Accept-Language")
		if lg != "" {
			// 语言设置不为空，就使用用户定义的语言
			lang.GetLang(lg)
		} else {
			// 否则就使用默认语言
			lang.GetLang(beego.AppConfig.String("lang"))
		}
		// 说明：当前的语言处理是根据，用户的请求信息决定的。你也可以根据具体的业务需求来决定，在什么位置设置语言。

		// token处理; token分为三个部分，其中两个部分分别是用户ID和token, 剩余的是干扰字符，当然也可以是简单的由ID和token组成。
		token := ctx.Request.Header.Get("X-Access-Token")
		var tok interface{}
		if err := dbs.RedisDB.GetJSON("test:user_1", &tok); err != nil {
			beego.Error(err)
		}

		// 这里是测试token匹配，实际需要根据用户ID读取redis里的token来进行匹配
		if token != "test" {
			// HTTP错误码 403 请求的资源不允许访问。就是说没有权限。
			ctx.ResponseWriter.WriteHeader(403)

			//定义返回对象
			var send controllers.SendMessage
			// 自定义错误码
			send.Code = 403
			// 自定义错误消息
			send.Message = lang.CurrLang.Routers.Filter.Token
			// 返回数据给用户
			if err := ctx.Output.JSON(&send, true, false); err != nil {
				// 打印错误信息。
				beego.Error(err)
			}
		}

		// IP处理
		ip := ctx.Request.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = ctx.Request.Header.Get("X-real-ip")
		}
		if ip == "" {
			ip = ctx.Input.IP()
		}
		if ip != "" {
			// 多语言设置：用户的IP地址:
			beego.Trace(lang.CurrLang.Routers.Filter.Ip.Info01, ip)
		} else {
			// 多语言设置：无法获取用户的IP地址:(
			beego.Trace(lang.CurrLang.Routers.Filter.Ip.Err01)
		}
	}
	// 启用过滤器
	// 第一个参数：表示拦截全部请求
	// 第二个参数：BeforeExec 找到路由之后，开始执行相应的 Controller 之前
	// 第三个参数：过滤器函数
	// 更多内容，请查看官方文档：https://beego.me/docs/mvc/controller/filter.md
	beego.InsertFilter("/*", beego.BeforeExec, check)
}
