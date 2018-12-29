package main

import (
	"github.com/astaxie/beego"
	"runtime"
	_ "testapi/conf"
	_ "testapi/dbs"
	_ "testapi/routers"
)

func main() {
	//指定使用多核，核心数为CPU的实际核心数量
	runtime.GOMAXPROCS(runtime.NumCPU())

	//beego应用启动
	beego.Run()
}
