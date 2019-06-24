package routers

import (
	"github.com/astaxie/beego"
	"goApiServer/controllers"
)

func init() {
	// api路由,一级
	ns := beego.NewNamespace("/test",
		// api路由,二级
		beego.NSNamespace("/api",
			//api路由,三级
			beego.NSNamespace(
				"/v1",
				// 用户登录
				beego.NSRouter("/login", &controllers.API{}, "post:Login"),
				// 用户登出
				beego.NSRouter("/logout", &controllers.API{}, "get:Logout"),
				// 查询用户, 分页
				beego.NSRouter("/users", &controllers.API{}, "get:SelectUsers"),
				// 查询用户
				beego.NSRouter("/user", &controllers.API{}, "get:SelectUser"),
				// 新增用户
				beego.NSRouter("/user", &controllers.API{}, "post:AddUser"),
				// 更新用户
				beego.NSRouter("/user", &controllers.API{}, "put:UpdateUser"),
				// 删除用户
				beego.NSRouter("/user", &controllers.API{}, "delete:DelUser"),
			),
		),
	)

	// websocket 路由
	ws := beego.NewNamespace("/ws",
		beego.NSRouter("/hi", &controllers.WebSocketController{}, "get:SayHi"),
	)

	// 注册自定义namespace
	beego.AddNamespace(ns, ws)

	// 执行过滤器
	RouterFilter()
}
