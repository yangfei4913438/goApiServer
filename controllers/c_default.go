package controllers

import (
	"github.com/astaxie/beego"
	"testapi/models"
)

func (api *API) Welcome() {

	//捕获get传参 /apy/api/v1/default?page_size=1&page_num=100
	page := api.GetString("page_size")
	num := api.GetString("page_num")
	beego.Notice(page, num)

	//自定义HTTP状态码
	api.Ctx.ResponseWriter.WriteHeader(200)

	//定义返回JSON
	api.Data["json"] = models.Default()

	//返回数据
	api.ServeJSON()
}
