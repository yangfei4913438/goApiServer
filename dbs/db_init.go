package dbs

func init() {
	initMysql()
	initRedis()
	initHttp()
}

/*
Http 请求范例

// 获取系统信息
res, err := dbs.Request.Get("/v1/system_info", nil, nil)
if err != nil {
	fmt.Println(err)
	return err
}
var data structs.ReceiveSystemInfo
_ = json.Unmarshal(res.Body, &data)

*/
