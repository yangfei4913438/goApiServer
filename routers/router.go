package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"testapi/controllers"
	"testapi/lang"
)

func init() {
	// api路由,一级
	ns := beego.NewNamespace("/test",
		// api路由,二级
		beego.NSNamespace("/api",
			//api路由,三级
			beego.NSNamespace(
				"/v1",
				//测试接口
				beego.NSRouter("/default", &controllers.API{}, "get:Welcome"),
				// 查询用户
				beego.NSRouter("/user", &controllers.API{}, "get:SelectUser"),
				// 新增用户
				beego.NSRouter("/user", &controllers.API{}, "put:AddUser"),
			),
		),
	)
	// 注册自定义namespace
	beego.AddNamespace(ns)

	// 过滤器函数，检查IP
	var checkIP = func(ctx *context.Context) {
		lg := ctx.Request.Header.Get("Accept-Language")
		res := lang.GetLang(lg)
		ip := ctx.Request.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = ctx.Request.Header.Get("X-real-ip")
		}
		if ip == "" {
			ip = ctx.Input.IP()
		}
		if ip != "" {
			beego.Trace(res.Routers.Filter.Ip.Info01, ip)
		} else {
			beego.Trace(res.Routers.Filter.Ip.Err01)
		}
	}
	// 启用过滤器
	beego.InsertFilter("/*", beego.BeforeExec, checkIP)
}
