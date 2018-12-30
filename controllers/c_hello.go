package controllers

func (api *API) Hello() {
	//自定义HTTP状态码
	api.Ctx.ResponseWriter.WriteHeader(200)

	//定义返回JSON
	api.Data["json"] = &SendMessage{200, "Hello World!"}

	//返回数据
	api.ServeJSON()
}
