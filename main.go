package main

import (
	"github.com/astaxie/beego"
	_ "goApiServer/conf"
	_ "goApiServer/dbs"
	_ "goApiServer/routers"
	"runtime"
)

func main() {
	//指定使用多核，核心数为CPU的实际核心数量
	runtime.GOMAXPROCS(runtime.NumCPU())

	//beego应用启动
	beego.Run()
}
