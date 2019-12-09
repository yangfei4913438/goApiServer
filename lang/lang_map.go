package lang

// JSON数据结构
type jsonData struct {
	ErrInfo ErrorInfo `json:"err_info"` // 错误信息
}

var CurrLang *jsonData

// 错误信息
type ErrorInfo struct {
	Err500  string `json:"err_500"`  // 服务器内部方法出错!
	Err1001 string `json:"err_1001"` // 用户名已经被注册!
	Err1002 string `json:"err_1002"` // 用户不存在!
	Err1003 string `json:"err_1003"` // 请求参数错误!
	Err1004 string `json:"err_1004"` // 请求的资源正忙，请稍后再试！
	Err1005 string `json:"err_1005"` // 身份验证失败!
	Err1006 string `json:"err_1006"` // 身份验证失败次数过多，请 5 分钟后再尝试登录。
	Err1007 string `json:"err_1007"` // 该账户因为登录错误太多被锁定，目前还处于安全锁定期间，请稍后再试！
	Err1008 string `json:"err_1008"` // 验证码不匹配!
}
