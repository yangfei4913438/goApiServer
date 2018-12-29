package controllers

import (
	"github.com/astaxie/beego"
)

type API struct {
	beego.Controller
}

type SendMessage struct {
	Errno  int    `json:"errno"`
	Errmsg string `json:"errmsg"`
}
