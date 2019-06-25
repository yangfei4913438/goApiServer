package structs

// 用户对象
type User struct {
	ID           int    `json:"id" db:"id"`                     // 用户ID, 主键
	Name         string `db:"name" json:"name"`                 // 用户昵称
	Password     string `db:"password" json:"password"`         // 登录密码
	Email        string `db:"email" json:"email"`               // 电子邮件
	Language     string `db:"language" json:"language"`         // 用户语言
	Role         int    `db:"role" json:"role"`                 // 用户角色编号 1 普通用户 9 管理员
	NoticeEnable int    `db:"noticeEnable" json:"noticeEnable"` // 是否开启邮件通知 -1 关闭 1 开启
	NoticeLevel  int    `db:"noticeLevel" json:"noticeLevel"`   // 邮件通知的级别
	CreateTime   int64  `db:"createTime" json:"createTime"`     // 用户创建时间
	CaptchaId    string `json:"captcha_id"`                     // 验证码的 ID
	ImgCode      string `json:"img_code"`                       // 用于接收验证码，数据库中不存在。
}

// 返回给用户的分页对象
type ReturnUsers struct {
	TotalNum int    `json:"total_num"`
	List     []User `json:"list"`
}

// 登录返回结构体
type LoginReturn struct {
	Token     string `json:"token"`      // 登陆成功后返回通信令牌
	CaptchaId string `json:"captcha_id"` // 登陆失败 3 次，返回验证码ID
	ImgCode   string `json:"img_code"`   // 验证码图片
	ErrCnt    int    `json:"err_cnt"`    // 登录失败的次数
	ErrCode   int    `json:"err_code"`   // 错误码
	ErrMsg    string `json:"err_msg"`    // 登录失败的错误消息
	UserInfo  *User  `json:"user_info"`  // 登陆成功后返回用户信息
}
