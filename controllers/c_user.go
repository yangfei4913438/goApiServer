package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"goApiServer/lang"
	"goApiServer/models"
	"goApiServer/structs"
	"goApiServer/tools"
)

func (api *API) AddUser() {
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
		code, err := models.AddUser(receive)
		if err != nil {
			if code == 0 {
				beego.Error("服务器内部方法出错:", err)
				// 自定义HTTP状态码
				api.Ctx.ResponseWriter.WriteHeader(500)
				// 自定义错误码和错误信息
				send = lang.RRSErrorInfo.ERR500
			} else {
				beego.Error("用户错误:", err)
				// 自定义HTTP状态码
				api.Ctx.ResponseWriter.WriteHeader(403)
				// 自定义错误码和错误信息
				send = &structs.SendMessage{Code: code, Message: err.Error()}
			}

		} else {
			//自定义HTTP状态码
			api.Ctx.ResponseWriter.WriteHeader(200)
			// 自定义返回信息
			send = &structs.SendMessage{Message: "用户添加成功!"}
		}
	}

	//定义返回JSON
	api.Data["json"] = send

	//返回数据
	res, _ := json.Marshal(send)
	beego.Info("返回给客户端的数据:", tools.BytesToString(res))
	api.ServeJSON()
}

// 查询单个用户的信息
func (api *API) SelectUser() {
	// 先设置语言
	api.GetLang()

	var send interface{}

	// 捕获get传参 /test/api/v1/user?id=1
	id := api.GetString("id", "0")
	// 调用查询方法
	user, code, err := models.GetUserInfo(id)
	if err != nil {
		if code == 0 {
			beego.Error("服务器内部方法出错:", err)
			// 自定义HTTP状态码
			api.Ctx.ResponseWriter.WriteHeader(500)
			// 自定义错误码和错误信息
			send = lang.RRSErrorInfo.ERR500
		} else {
			beego.Error("用户错误:", err)
			// 自定义HTTP状态码
			api.Ctx.ResponseWriter.WriteHeader(403)
			// 自定义错误码和错误信息
			send = &structs.SendMessage{Code: code, Message: err.Error()}
		}

	} else {
		//自定义HTTP状态码
		api.Ctx.ResponseWriter.WriteHeader(200)
		// 返回的内容
		send = user
	}

	//定义返回JSON
	api.Data["json"] = send

	//返回数据
	res, _ := json.Marshal(send)
	//多语言打印：返回给客户端的数据:
	beego.Info("返回给客户端的数据:", string(res))
	api.ServeJSON()
}

func (api *API) SelectUsers() {
	// 先设置语言
	api.GetLang()

	var send interface{}

	// 捕获get传参 /server/api/v1/user?page_num=1&page_size=10

	// 获取分页的页数
	number := api.GetString("page_num")
	// 获取分页的每页数量
	size := api.GetString("page_size")

	// 调用查询方法
	users, err := models.SelectUsers(number, size)
	if err != nil {
		beego.Error("服务器内部方法出错:", err)
		// 自定义HTTP状态码
		api.Ctx.ResponseWriter.WriteHeader(500)
		// 自定义错误码和错误信息
		send = lang.RRSErrorInfo.ERR500
	} else {
		//自定义HTTP状态码
		api.Ctx.ResponseWriter.WriteHeader(200)
		// 返回的内容
		send = users
	}

	//定义返回JSON
	api.Data["json"] = send

	//返回数据
	res, _ := json.Marshal(send)
	//多语言打印：返回给客户端的数据:
	beego.Info("返回给客户端的数据:", string(res))
	api.ServeJSON()
}

func (api *API) UpdateUser() {
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
		code, err2 := models.UpdateUser(receive)
		if err2 != nil {
			if code == 0 {
				beego.Error("服务器内部方法出错:", err2)
				// 自定义HTTP状态码
				api.Ctx.ResponseWriter.WriteHeader(500)
				// 自定义错误码和错误信息
				send = lang.RRSErrorInfo.ERR500
			} else {
				beego.Error("用户错误:", err2)
				// 自定义HTTP状态码
				api.Ctx.ResponseWriter.WriteHeader(403)
				// 自定义错误码和错误信息
				send = &structs.SendMessage{Code: code, Message: err2.Error()}
			}
		} else {
			//自定义HTTP状态码
			api.Ctx.ResponseWriter.WriteHeader(200)
			// 自定义返回信息
			send = &structs.SendMessage{Message: "用户更新成功!"}
		}
	}

	//定义返回JSON
	api.Data["json"] = send

	//返回数据
	res, _ := json.Marshal(send)
	beego.Info("返回给客户端的数据:", string(res))
	api.ServeJSON()
}

func (api *API) DelUser() {
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
		code, err2 := models.DelUser(receive.ID)
		if err2 != nil {
			if code == 0 {
				beego.Error("服务器内部方法出错:", err2)
				// 自定义HTTP状态码
				api.Ctx.ResponseWriter.WriteHeader(500)
				// 自定义错误码和错误信息
				send = lang.RRSErrorInfo.ERR500
			} else {
				beego.Error("用户错误:", err2)
				// 自定义HTTP状态码
				api.Ctx.ResponseWriter.WriteHeader(403)
				// 自定义错误码和错误信息
				send = &structs.SendMessage{Code: code, Message: err2.Error()}
			}
		} else {
			//自定义HTTP状态码
			api.Ctx.ResponseWriter.WriteHeader(200)
			// 自定义返回信息
			send = &structs.SendMessage{Message: "用户删除成功!"}
		}
	}

	//定义返回JSON
	api.Data["json"] = send

	//返回数据
	res, _ := json.Marshal(send)
	beego.Info("返回给客户端的数据:", string(res))
	api.ServeJSON()
}
