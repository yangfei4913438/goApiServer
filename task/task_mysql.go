package task

import (
	"fmt"
	"github.com/astaxie/beego"
	"rs-server/dbs"
	"time"
)

func delOpLogs() error {
	beego.Notice("开始维护操作日志表。。。")
	defer beego.Notice("操作日志表维护完成。")

	// 读取数据库的时间配置
	// 组合查询sql
	sql := "select op_log_expired from system_config"
	beego.Debug("[sql]: " + sql)

	// 过期时间变量
	var expired int64
	err := dbs.MysqlDB.Get(&expired, sql)
	if err != nil {
		beego.Error(err)
		return err
	}
	beego.Debug("当前的日志保留时间为:", expired, "天。")

	var expiredTime int64
	// 计算指定的日期
	expiredTime = time.Now().Unix() - expired*24*60*60

	// 拼接sql
	sql = fmt.Sprintf("delete from user_log where created_at < %v", expiredTime)

	// 打印日志
	beego.Debug("[sql]: " + sql)

	// 执行sql
	_, err = dbs.MysqlDB.Exec(sql)
	if err != nil {
		// 打印错误日志
		beego.Error("执行删除 sql 出错:", err)
		// 返回错误信息
		return err
	}

	return nil
}
