package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"testapi/models"
)

func (api *API) SelectUser() {
	var send interface{}

	// 捕获get传参 /test/api/v1/user?id=1
	id, err := api.GetInt64("id")
	if err != nil {
		beego.Error("服务器获取参数出错:", err)
		// 自定义HTTP状态码
		api.Ctx.ResponseWriter.WriteHeader(500)
		// 自定义错误码和错误信息
		send = &SendMessage{101, "服务器获取参数出错!"}
	} else {
		// 调用查询方法
		user, err2 := models.SelectUser(id)
		if err2 != nil {
			beego.Error("服务器内部方法出错:", err2)
			// 自定义HTTP状态码
			api.Ctx.ResponseWriter.WriteHeader(500)
			// 自定义错误码和错误信息
			send = &SendMessage{102, "服务器内部方法出错!"}
		} else {
			// 返回的内容
			send = user
			//自定义HTTP状态码
			api.Ctx.ResponseWriter.WriteHeader(200)
		}
	}

	//定义返回JSON
	api.Data["json"] = send

	//返回数据
	res, _ := json.Marshal(send)
	beego.Trace("返回给客户端的数据:", string(res))
	api.ServeJSON()
}

func (api *API) AddUser() {
	// 接收结构体
	var receive *models.ReceiveUser
	// 消息对象
	var send interface{}

	if err := json.Unmarshal(api.Ctx.Input.RequestBody, &receive); err != nil {
		beego.Error("json反序列化出错: " + err.Error())
		//自定义HTTP状态码
		api.Ctx.ResponseWriter.WriteHeader(500)
		send = &SendMessage{103, "json反序列化出错"}
	} else {
		// 调用新增用户方法
		err2 := models.AddUser(receive)
		if err2 != nil {
			beego.Error("服务器内部方法出错:", err2)
			// 自定义HTTP状态码
			api.Ctx.ResponseWriter.WriteHeader(500)
			// 自定义错误码和错误信息
			send = &SendMessage{102, "服务器内部方法出错!"}
		} else {
			//自定义HTTP状态码
			api.Ctx.ResponseWriter.WriteHeader(200)
			// 自定义返回信息
			send = &SendMessage{0, "用户添加成功!"}
		}
	}

	//定义返回JSON
	api.Data["json"] = send

	//返回数据
	api.ServeJSON()
}
