package controllers

import (
	"github.com/astaxie/beego"
)

type API struct {
	beego.Controller
}

type SendMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
