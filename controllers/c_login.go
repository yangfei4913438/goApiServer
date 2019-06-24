package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"goApiServer/lang"
	"goApiServer/models"
	"goApiServer/structs"
)

func (api *API) Login() {
	// 先设置语言
	api.GetLang()

	// 接收结构体
	var receive *structs.User
	// 消息对象
	var send interface{}

	if err := json.Unmarshal(api.Ctx.Input.RequestBody, &receive); err != nil {
		beego.Error("json反序列化出错: " + err.Error())
		//自定义HTTP状态码
		api.Ctx.ResponseWriter.WriteHeader(500)
		send = lang.RRSErrorInfo.ERR500
	} else {
		// 调用新增用户方法
		res := models.Login(receive)
		if res.ErrMsg != "" {
			if res.ErrCode == 0 {
				beego.Error("服务器内部方法出错:", res.ErrMsg)
				// 自定义HTTP状态码
				api.Ctx.ResponseWriter.WriteHeader(500)
				// 自定义错误码和错误信息
				send = lang.RRSErrorInfo.ERR500
			} else {
				beego.Error("用户错误:", res.ErrMsg)
				// 自定义HTTP状态码
				api.Ctx.ResponseWriter.WriteHeader(403)
				// 自定义错误码和错误信息
				send = res
			}
		} else {
			//自定义HTTP状态码
			api.Ctx.ResponseWriter.WriteHeader(200)
			// 自定义返回信息
			send = res
		}
	}

	//定义返回JSON
	api.Data["json"] = send

	//返回数据
	res, _ := json.Marshal(send)
	beego.Info("返回给客户端的数据:", string(res))
	api.ServeJSON()
}

func (api *API) Logout() {
	// 先设置语言
	api.GetLang()

	// 接收结构体
	token := api.Ctx.Request.Header.Get("X-Access-Token")
	// 消息对象
	var send interface{}

	// 调用新增用户方法
	err := models.Logout(token)
	if err != nil {
		// 因为只要进入到这里，证明用户传递的 token没有问题，只要出错，肯定是服务器内部函数错误。
		beego.Error("服务器内部方法出错:", err)
		// 自定义HTTP状态码
		api.Ctx.ResponseWriter.WriteHeader(500)
		// 自定义错误码和错误信息
		send = lang.RRSErrorInfo.ERR500
	} else {
		//自定义HTTP状态码
		api.Ctx.ResponseWriter.WriteHeader(200)
		// 自定义返回信息
		send = &structs.SendMessage{Message: "用户登出成功!"}
	}

	//定义返回JSON
	api.Data["json"] = send

	//返回数据
	res, _ := json.Marshal(send)
	beego.Info("返回给客户端的数据:", string(res))
	api.ServeJSON()
}
