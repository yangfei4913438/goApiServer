package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"goApiServer/dbs"
	"goApiServer/lang"
	"goApiServer/structs"
	"goApiServer/tools"
	"strings"
)

func RouterFilter() {
	// 过滤器函数，检查操作
	var check = func(ctx *context.Context) {
		// 设置语言
		var userLang string
		// 语言处理
		lg := ctx.Request.Header.Get("Accept-Language")
		if lg != "" {
			// 语言设置不为空，就使用用户定义的语言
			userLang = lg
		} else {
			// 否则就使用默认语言
			userLang = beego.AppConfig.String("lang")
		}
		// 因为还没有到 API 的流程，所以这里要单独设置一下语言。
		lang.SetLang(userLang)

		// 如果是登陆请求，不处理 token
		var pathKey []string
		pathKey = strings.Split(ctx.Request.URL.Path, "/")
		// 切片的最后一个字符串进行匹配, 因为有可能是多个版本的登陆接口，所以只能匹配最后一节字符串
		if pathKey[len(pathKey)-1] != "login" {
			// token处理; token分为三个部分，其中两个部分分别是用户ID和token, 剩余的是干扰字符，当然也可以是简单的由ID和token组成。
			token := ctx.Request.Header.Get("X-Access-Token")
			// 如果为空，表示没取到，那么就表示这个请求可能是 websocket 请求，我们到 url里面查询一下看看
			if token == "" {
				// 将 url 中查询到的值，赋值给变量
				token = ctx.Input.Query("token")
			}
			// 说明：当前的规则下，只有这两种方式可以正确的获取到 token，其他情况不予考虑，如果还没有获取到，那么表示用户没有权限。

			// 查询 token 是否存在缓存中, token 在存储的时候，会添加 user_token: 前缀。避免用户随意传一个已经存在的值过来。
			tokenKey := "user_token:" + token
			ok, _ := dbs.RedisDB.Exists(tokenKey)

			// 如果没有查询到当前 token，那么就返回错误信息
			if !ok {
				// HTTP错误码 401 没有提供认证信息。请求的时候没有带上 Token，或者 token 错误等。
				ctx.ResponseWriter.WriteHeader(401)

				//定义返回对象
				var send structs.SendMessage
				// 自定义错误码
				send.Code = 1001
				// 自定义错误消息
				send.Message = lang.CurrLang.ErrInfo.Err1005
				// 打印日志
				beego.Info("返回给客户端的数据:", tools.ReturnJson(send))
				// 返回数据给用户
				if err := ctx.Output.JSON(&send, true, false); err != nil {
					// 打印错误信息。
					beego.Error(err)
					return
				}
			} else {
				//	如果查询到了 token, 那么就重置 token 的有效期，默认为 2 小时
				_ = dbs.RedisDB.SetLife(tokenKey, tools.OneHour*2)
			}
		}
	}
	// 启用过滤器
	// 第一个参数：表示拦截全部请求
	// 第二个参数：BeforeExec 找到路由之后，开始执行相应的 Controller 之前
	// 第三个参数：过滤器函数
	// 更多内容，请查看官方文档：https://beego.me/docs/mvc/controller/filter.md
	beego.InsertFilter("/*", beego.BeforeExec, check)
}
